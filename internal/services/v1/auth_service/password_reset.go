package auth_service

import (
	"bytes"
	"context"
	"errors"
	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/pkg/helpers"
	"github.com/abdivasiyev/project_template/pkg/mailer"
	"github.com/abdivasiyev/project_template/pkg/validator"
	"go.uber.org/zap"
	"html/template"
	"time"
)

const resetPasswordEmailTemplate = `
<!doctype html>
<html lang="en-US">
<head>
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type" />
    <title>Reset Password Email</title>
    <meta name="description" content="Reset Password Email Template.">
    <style type="text/css">
        a:hover {text-decoration: underline !important;}
    </style>
</head>
<body marginheight="0" topmargin="0" marginwidth="0" style="margin: 0px; background-color: #f2f3f8;" leftmargin="0">
    <table cellspacing="0" border="0" cellpadding="0" width="100%" bgcolor="#f2f3f8"
        style="@import url(https://fonts.googleapis.com/css?family=Rubik:300,400,500,700|Open+Sans:300,400,600,700); font-family: 'Open Sans', sans-serif;">
        <tr>
            <td>
                <table style="background-color: #f2f3f8; max-width:670px;  margin:0 auto;" width="100%" border="0"
                    align="center" cellpadding="0" cellspacing="0">
                    <tr>
                        <td style="height:80px;">&nbsp;</td>
                    </tr>
                    <tr>
                        <td style="text-align:center;">
                          <a href="https://example.com" title="logo" target="_blank">
                            <img width="300" src="https://example.com/uploads/logo.png" title="logo" alt="logo">
                          </a>
                        </td>
                    </tr>
                    <tr>
                        <td style="height:20px;">&nbsp;</td>
                    </tr>
                    <tr>
                        <td>
                            <table width="95%" border="0" align="center" cellpadding="0" cellspacing="0"
                                style="max-width:670px;background:#fff; border-radius:3px; text-align:center;-webkit-box-shadow:0 6px 18px 0 rgba(0,0,0,.06);-moz-box-shadow:0 6px 18px 0 rgba(0,0,0,.06);box-shadow:0 6px 18px 0 rgba(0,0,0,.06);">
                                <tr>
                                    <td style="height:40px;">&nbsp;</td>
                                </tr>
                                <tr>
                                    <td style="padding:0 35px;">
                                        <h1 style="color:#1e1e2d; font-weight:500; margin:0;font-size:32px;font-family:'Rubik',sans-serif;">You have requested to reset your password</h1>
                                        <span
                                            style="display:inline-block; vertical-align:middle; margin:29px 0 26px; border-bottom:1px solid #cecece; width:100px;"></span>
                                        <p style="color:#455056; font-size:15px;line-height:24px; margin:0;">
                                            We cannot simply send you your old password. A temporary password to reset your
                                            password has been generated for you. To reset your password, enter the
                                          following code to <b>Verification Code</b> input in application or web site.
                                        </p>
                                        <a href="javascript:void(0);" style="background:#20e277;text-decoration:none !important; font-weight:500; margin-top:35px; color:#fff;text-transform:uppercase; font-size:24px;padding:10px 24px;display:inline-block;border-radius:50px;">
											{{.PasswordResetCode}}
										</a>
                                    </td>
                                </tr>
                                <tr>
                                    <td style="height:40px;">&nbsp;</td>
                                </tr>
                            </table>
                        </td>
                    <tr>
                        <td style="height:20px;">&nbsp;</td>
                    </tr>
                    <tr>
                        <td style="text-align:center;">
                            <p style="font-size:14px; color:rgba(69, 80, 86, 0.7411764705882353); line-height:18px; margin:0 0 0;">&copy; <strong>example.com</strong></p>
                        </td>
                    </tr>
                    <tr>
                        <td style="height:80px;">&nbsp;</td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`

func (s *service) ResetPassword(ctx context.Context, req models.ResetPasswordRequest) (models.SuccessResponse, error) {
	var (
		key                = "auth:reset_password:" + req.Email
		resetPasswordCache models.ResetPasswordCache
	)

	if err := s.cache.GetObj(ctx, key, &resetPasswordCache); err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get password reset from cache", zap.Error(err))
			return models.SuccessResponse{}, err
		}
	}

	// if user id empty find user and send email
	if helpers.IsEmpty(resetPasswordCache.ResetCode) {
		err := s.sendResetCode(ctx, key, req.Email)
		if err != nil {
			s.sentry.HandleError(err)
			s.log.Error("could not get password reset from cache", zap.Error(err))
			return models.SuccessResponse{}, err
		}

		return models.SuccessResponse{Ok: true}, nil
	}

	if !helpers.IsEmpty(resetPasswordCache.ResetCode) && !helpers.IsEmpty(req.ResetCode) {
		err := s.verifyResetCode(ctx, key, resetPasswordCache, req.ResetCode)
		if err != nil {
			s.sentry.HandleError(err)
			s.log.Error("could not get password reset from cache", zap.Error(err))
			return models.SuccessResponse{}, err
		}

		return models.SuccessResponse{Ok: true}, nil
	}

	if resetPasswordCache.Verified && !helpers.IsEmpty(req.Password) {
		err := s.updateUserPassword(ctx, key, resetPasswordCache.UserID, req.Password)
		if err != nil {
			s.sentry.HandleError(err)
			s.log.Error("could not get password reset from cache", zap.Error(err))
			return models.SuccessResponse{}, err
		}

		return models.SuccessResponse{Ok: true}, nil
	}

	// update user password if password is not empty

	return models.SuccessResponse{}, validator.NewValidationError("email", "something went wrong")
}

func (s *service) updateUserPassword(ctx context.Context, cacheKey, userID, password string) error {
	user, err := s.userRepository.Get(ctx, userID)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.log.Error("could not get user", zap.Error(err))
			s.sentry.HandleError(err)
		}
		return err
	}

	passwordHash, err := s.security.GenerateHash(password)
	if err != nil {
		s.log.Error("could not generate password hash", zap.Error(err))
		s.sentry.HandleError(err)
		return err
	}

	err = s.userRepository.Update(ctx, models.UpdateUserRequest{
		ID:          user.ID,
		CompanyID:   user.Company.ID,
		RoleID:      user.Role.ID,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		NewPassword: passwordHash,
	})
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.log.Error("could not update user", zap.Error(err))
			s.sentry.HandleError(err)
		}
		return err
	}

	if err = s.cache.Delete(ctx, cacheKey); err != nil {
		s.log.Error("could not get user", zap.Error(err))
		s.sentry.HandleError(err)
	}

	return nil
}

func (s *service) verifyResetCode(ctx context.Context, cacheKey string, resetPasswordCache models.ResetPasswordCache, resetCode string) error {
	if resetPasswordCache.ResetCode != resetCode {
		return validator.NewValidationError("reset_code", "reset code is not valid")
	}

	resetPasswordCache.Verified = true
	if err := s.cache.SetObj(ctx, cacheKey, resetPasswordCache, 10*time.Minute); err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not set reset password cache", zap.Error(err))
		return err
	}

	return nil
}

func (s *service) sendResetCode(ctx context.Context, cacheKey, email string) error {
	var resetPasswordCache models.ResetPasswordCache
	user, err := s.userRepository.GetByUsername(ctx, email)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not get password reset from cache", zap.Error(err))
		return err
	}
	resetPasswordCache.UserID = user.ID

	resetCode, err := helpers.RandomNumber(10_000)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not generate random number", zap.Error(err))
		return err
	}

	tmp := template.New("password_reset")

	tmp, err = tmp.Parse(resetPasswordEmailTemplate)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not parse email template", zap.Error(err))
		return err
	}

	buf := &bytes.Buffer{}

	err = tmp.Execute(buf, struct {
		PasswordResetCode string
	}{
		PasswordResetCode: resetCode,
	})
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not execute template", zap.Error(err))
		return err
	}

	if err = s.mailer.Send(mailer.Mail{
		To:      []string{email},
		Subject: "Password ResetPassword Code",
		Body:    buf.String(),
	}); err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not send email", zap.Error(err))
		return err
	}
	resetPasswordCache.ResetCode = resetCode
	resetPasswordCache.Email = email

	if err = s.cache.SetObj(ctx, cacheKey, resetPasswordCache, 10*time.Minute); err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not set reset password cache", zap.Error(err))
		return err
	}

	return nil
}
