package http

import (
	"context"
	"strconv"

	"github.com/hdkef/be-assignment/services/account/domain/entity"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func (h *HttpHandler) SuperTokenSignUp(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
	originalSignUpPOST := *originalImplementation.SignUpPOST

	(*originalImplementation.SignUpPOST) = func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {

		// build dto
		dto := entity.UserSignUpDto{}
		for _, field := range formFields {
			switch field.ID {
			case "email":
				dto.Email = field.Value
			case "name":
				dto.Name = field.Value
			case "dateOfBirth":
				err := dto.SetDOB(field.Value)
				if err != nil {
					return epmodels.SignUpPOSTResponse{}, err
				}
			case "job":
				dto.Job = field.Value
			case "address":
				dto.Address = field.Value
			case "district":
				dto.District = field.Value
			case "city":
				dto.City = field.Value
			case "province":
				dto.Province = field.Value
			case "country":
				dto.Country = field.Value
			case "accCurrency":
				dto.FirstAccountCurrency = field.Value
			case "accDesc":
				dto.FirstAccountDesc = field.Value
			case "zip":
				val, _ := strconv.Atoi(field.Value)
				dto.ZIP = uint32(val)
			}
		}

		// pre API logic...

		resp, err := originalSignUpPOST(formFields, tenantId, options, userContext)
		if err != nil {
			return epmodels.SignUpPOSTResponse{}, err
		}

		if resp.OK != nil {
			err = dto.SetUserID(resp.OK.User.ID)
			if err != nil {
				// TODO implement rollback supertoken
				return epmodels.SignUpPOSTResponse{}, err
			}

			// execute sign up usecase
			err = h.UserUsecase.SignUp(context.Background(), &dto)
			if err != nil {
				// TODO implement rollback supertoken
				return epmodels.SignUpPOSTResponse{}, err
			}

		}

		return resp, nil
	}

	return originalImplementation
}
