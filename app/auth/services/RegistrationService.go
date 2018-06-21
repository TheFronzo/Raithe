package services

import (
	"github.com/labstack/echo"
	"encoding/json"
	"github.com/catmullet/Raithe/app/auth/model"
	"fmt"
	"crypto/rand"
	"crypto/rsa"
	"github.com/catmullet/Raithe/app/utils"
)

var (
	RegisteredAgents []model.SecurityToken
)

func getAgents() model.Agents {
	return utils.GetAgentsFromList()
}

func RegisterAsAgent(ctx echo.Context) error {
	reg := model.Register{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&reg)

	if err != nil {
		fmt.Println(err)
	}

	agents := getAgents()

	if isAlreadyRegistered(reg.AgentName) {
		return ctx.JSON(200, model.RegisterResponse{Success:false, Message:"Agent is already Registered"})
	}

	for _, val := range agents.Agents {
		if val == reg.AgentName {
			token, _ := GeneratePrivateKey()

			secToken := model.SecurityToken{AgentName:reg.AgentName,Token:token}
			RegisteredAgents = append(RegisteredAgents, secToken)

			return ctx.JSON(200, model.RegisterResponse{Success:true, SecurityToken:secToken})
		}
	}

	return ctx.JSON(200, model.RegisterResponse{Success:false,Message:"Unrecognized Agent"})
}


func isAlreadyRegistered(agentName string) bool {
	for _, val := range RegisteredAgents {
		if val.AgentName == agentName {
			return true
		}
	}

	return false
}

func GeneratePrivateKey() (string, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", privateKey.D.Bytes()), nil
}

func IsAgentRegistered(token model.SecurityToken) bool {

	for _, val := range RegisteredAgents {
		if val.Token == token.Token && val.AgentName == token.AgentName {
			return true
		}
	}
	return false
}

func InvalidateTokens(ctx echo.Context) error {
	inv := model.InvalidateTokens{}
	err := ctx.Bind(&inv)

	if err != nil {
		return err
	}

	if !IsAgentRegistered(inv.Token){
		return ctx.JSON(403, model.ValidateResponse{Success:false, Message:"Security Token Not Recognized"})
	}
	RegisteredAgents = []model.SecurityToken{}
	return ctx.JSON(200, "Invalidated Tokens")
}

func DumpTokens(ctx echo.Context) error {
	inv := model.InvalidateTokens{}
	err := ctx.Bind(&inv)

	if err != nil {
		return err
	}

	if !IsAgentRegistered(inv.Token){
		return ctx.JSON(403, model.ValidateResponse{Success:false, Message:"Security Token Not Recognized"})
	}

	for _, val := range RegisteredAgents {
		fmt.Println(val)
	}
	return ctx.JSON(200, "Tokens Have been dumped to logs")
}
