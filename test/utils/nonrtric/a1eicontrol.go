package nonrtric

import "github.com/labstack/echo/v4"

/*
	EI Implementations of Controller interface
*/

// EI job identifiers
// (GET /A1-EI/v1/eijobs)
func (c *controller) A1EIGetEiJobIdsUsingGET(ctx echo.Context, eiTypeId string) error {
	return nil
}

// Individual EI job
// (DELETE /A1-EI/v1/eijobs/{eiJobId})
func (c *controller) A1EIDeleteIndividualEiJobUsingDELETE(ctx echo.Context, eiJobId string) error {
	return nil
}

// Individual EI job
// (GET /A1-EI/v1/eijobs/{eiJobId})
func (c *controller) A1EIGetIndividualEiJobUsingGET(ctx echo.Context, eiJobId string) error {
	return nil
}

// Individual EI job
// (PUT /A1-EI/v1/eijobs/{eiJobId})
func (c *controller) A1EIPutIndividualEiJobUsingPUT(ctx echo.Context, eiJobId string) error {
	return nil
}

// EI job status
// (GET /A1-EI/v1/eijobs/{eiJobId}/status)
func (c *controller) A1EIGetEiJobStatusUsingGET(ctx echo.Context, eiJobId string) error { return nil }

// EI type identifiers
// (GET /A1-EI/v1/eitypes)
func (c *controller) A1EIGetEiTypeIdentifiersUsingGET(ctx echo.Context) error { return nil }

// Individual EI type
// (GET /A1-EI/v1/eitypes/{eiTypeId})
func (c *controller) A1EIGetEiTypeUsingGET(ctx echo.Context, eiTypeId string) error { return nil }
