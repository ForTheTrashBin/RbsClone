package rest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/ForTheTrashBin/RbsClone/internal/rbsdb"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

//-----------------------------------------------------------------------------
//
//	GET		Get List 			"/custodian"
//	GET		Get by ID			"/custodian/id/{id}"
//	GET		Get by Shortcode	"/custodian/shortcode/{shortcode}"
//	POST	Create				"/custodian"
//	DELETE	Delete				"/custodian/{id}"
//	PUT		Update				"/custodian/{id}"
//
//-----------------------------------------------------------------------------

type CustodianId struct {
	ID uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type CustodianIdPath struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type CustodianShortcode struct {
	Shortcode string `json:"shortcode" minLength:"1" maxLength:"5" doc:"A unique short name for this data"`
}

type CustodianShortcodePath struct {
	Shortcode string `path:"shortcode" minLength:"1" maxLength:"5" doc:"A unique short name for this data"`
}

type CustodianName struct {
	Name string `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
}

type CustodianFlags struct {
	Flags int16 `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
}

type CustodianIdCountry struct {
	Idcountry uuid.UUID `json:"idcountry" format:"uuid" doc:"A reference to a country, where the custodion is in"`
}

type CustodianDepotNo struct {
	Depotno *string `json:"depotno,omitempty" maxLength:"10" doc:"This dopot numer assocciated with this custodian"`
}

type Exchange2CustodianIdExchange struct {
	Idexchange uuid.UUID `json:"idexchange" format:"uuid" doc:"This is the unique identifier of this data"`
}

type Exchange2CustodianIdExchangeDefault struct {
	Idexchangedefault *uuid.UUID `json:"idexchangedefault,omitempty" format:"uuid" doc:"The id of the default-exchange"`
}

type Exchange2CustodianValue1 struct {
	Value1 int16 `json:"value1" format:"int16" minimum:"0" maximum:"255" doc:"This is the first special payload for testing"`
}

type Exchange2CustodianValue2 struct {
	Value2 int16 `json:"value2" format:"int16" minimum:"0" maximum:"255" doc:"This is the second special payload for testing"`
}

//-----------------------------------------------------------------------------

type CustodianListItem struct {
	CustodianId
	CustodianShortcode
	CustodianName
}

type Exchange2CustodianListItem struct {
	Exchange2CustodianIdExchange
	Exchange2CustodianValue1
	Exchange2CustodianValue2
}

type Exchange2CustodianList struct {
	Exchanges []Exchange2CustodianListItem `json:"exchanges"`
}

type CustodianNoPK struct {
	CustodianShortcode
	CustodianName
	CustodianFlags
	CustodianIdCountry
	CustodianDepotNo
	Exchange2CustodianIdExchangeDefault
	Exchange2CustodianList
}

type Custodian struct {
	CustodianId
	CustodianShortcode
	CustodianName
	CustodianFlags
	CustodianIdCountry
	CustodianDepotNo
	Exchange2CustodianIdExchangeDefault
	Exchange2CustodianList
}

//-----------------------------------------------------------------------------

type CustodianRequestId struct {
	CustodianIdPath
}

type CustodianRequestShortcode struct {
	CustodianShortcodePath
}

type CustodianRequestCreate struct {
	Body CustodianNoPK
}

type CustodianRequestUpdate struct {
	CustodianIdPath
	Body CustodianNoPK
}

//-----------------------------------------------------------------------------

type CustodianResponseList struct {
	Body []CustodianListItem
}

type CustodianResponseCreate struct {
	Header LocationHeader
	Body   Custodian
}

type CustodianResponse struct {
	Body Custodian
}

//-----------------------------------------------------------------------------

func (rs *RestServer) getExchange2CustodianDefaultList(ctx context.Context, IdCustodian uuid.UUID) (*Exchange2CustodianList, *uuid.UUID, error) {

	//-------------------------------------------------------------------------
	// Read a slice of records from database (could be zero)
	//-------------------------------------------------------------------------

	dbExchange2Custodian, err := rs.dbQueries.GetExchange2CustodianByIdcustodian(ctx, IdCustodian)

	if err != nil {

		if !errors.Is(err, sql.ErrNoRows) {

			return nil, nil, err
		}
	}

	//-------------------------------------------------------------------------
	// If there are exchanges, the first is the default exchange
	//-------------------------------------------------------------------------

	var idexchangedefault uuid.UUID = uuid.Nil

	if len(dbExchange2Custodian) > 0 {

		idexchangedefault = dbExchange2Custodian[0].Idexchange
	}

	//-------------------------------------------------------------------------
	// Create the exchange-list and fill it with results from database
	//-------------------------------------------------------------------------

	exchange2CustodianList := Exchange2CustodianList{

		Exchanges: make([]Exchange2CustodianListItem, len(dbExchange2Custodian)),
	}

	for idx, dbData := range dbExchange2Custodian {

		exchange2CustodianList.Exchanges[idx] = Exchange2CustodianListItem{
			Exchange2CustodianIdExchange: Exchange2CustodianIdExchange{Idexchange: dbData.Idexchange},
			Exchange2CustodianValue1:     Exchange2CustodianValue1{Value1: dbData.Value1},
			Exchange2CustodianValue2:     Exchange2CustodianValue2{Value2: dbData.Value2},
		}
	}

	//-------------------------------------------------------------------------

	return &exchange2CustodianList, &idexchangedefault, nil
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerCustodianRoutes() {

	group := huma.NewGroup(rs.api, "/custodian")

	group.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {

		op.Tags = append(op.Tags, "Custodian")

		next(op)
	})

	//-------------------------------------------------------------------------

	huma.Get(group, "", func(ctx context.Context, request *struct{}) (*CustodianResponseList, error) {

		OperationId := getOperationIdFromContext(ctx)

		dbSlice, err := rs.dbQueries.GetCustodians(ctx)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		result := make([]CustodianListItem, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APICustodianListItem(dbData)
		}

		return &CustodianResponseList{Body: result}, nil

	}, describeEndpoint("getCustodians", "Get a list of all custodians"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/id/{id}", func(ctx context.Context, request *CustodianRequestId) (*CustodianResponse, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		//---------------------------------------------------------------------
		// Read a single record from database
		//---------------------------------------------------------------------

		dbCustodian, err := rs.dbQueries.GetCustodianById(ctx, request.ID)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, OperationId, "error", err)

				return nil, mapDBError(err)
			}
		}

		//---------------------------------------------------------------------
		// Read the Exchange2CustodianDefaultList
		//---------------------------------------------------------------------

		E2C_List, E2C_IdExchangeDefault, err := rs.getExchange2CustodianDefaultList(ctx, dbCustodian.ID)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		//---------------------------------------------------------------------

		if E2C_List == nil { // Should never be nil at this point

			rs.logger.ErrorContext(ctx, OperationId, "error", fmt.Errorf("E2C_List is nil"))

			return nil, huma.Error500InternalServerError("Internal server error")
		}

		if E2C_IdExchangeDefault == nil { // Should never be nil at this point, but uuid.Nil at least

			rs.logger.ErrorContext(ctx, OperationId, "error", fmt.Errorf("E2C_IdExchangeDefault is nil"))

			return nil, huma.Error500InternalServerError("Internal server error")
		}

		//---------------------------------------------------------------------

		return &CustodianResponse{Body: mapDB2APICustodian(dbCustodian, *E2C_List, *E2C_IdExchangeDefault)}, nil

	}, describeEndpoint("getCustodianById", "Get a single custodian based on the id supplied"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/shortcode/{shortcode}", func(ctx context.Context, request *CustodianRequestShortcode) (*CustodianResponse, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		//---------------------------------------------------------------------
		// Read a single record from database
		//---------------------------------------------------------------------

		dbCustodian, err := rs.dbQueries.GetCustodianByShortcode(ctx, request.Shortcode)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, OperationId, "error", err)

				return nil, mapDBError(err)
			}
		}

		//---------------------------------------------------------------------
		// Read the Exchange2CustodianDefaultList
		//---------------------------------------------------------------------

		E2C_List, E2C_IdExchangeDefault, err := rs.getExchange2CustodianDefaultList(ctx, dbCustodian.ID)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		//---------------------------------------------------------------------

		if E2C_List == nil { // Should never be nil at this point

			rs.logger.ErrorContext(ctx, OperationId, "error", fmt.Errorf("E2C_List is nil"))

			return nil, huma.Error500InternalServerError("Internal server error")
		}

		if E2C_IdExchangeDefault == nil { // Should never be nil at this point, but uuid.Nil at least

			rs.logger.ErrorContext(ctx, OperationId, "error", fmt.Errorf("E2C_IdExchangeDefault is nil"))

			return nil, huma.Error500InternalServerError("Internal server error")
		}

		//---------------------------------------------------------------------

		return &CustodianResponse{Body: mapDB2APICustodian(dbCustodian, *E2C_List, *E2C_IdExchangeDefault)}, nil

	}, describeEndpoint("getCustodianByShortcode", "Get a single custodian based on the shortcode supplied"))

	//-------------------------------------------------------------------------

	huma.Post(group, "", func(ctx context.Context, request *CustodianRequestCreate) (*CustodianResponseCreate, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		//---------------------------------------------------------------------
		// Initialize transaction
		//---------------------------------------------------------------------

		tx, err := rs.dbPool.Begin(ctx)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		defer tx.Rollback(ctx)

		queriesWithTx := rs.dbQueries.WithTx(tx)

		//---------------------------------------------------------------------
		// Insert custodian into database
		//---------------------------------------------------------------------

		insertParams := rbsdb.InsertCustodianParams{

			Shortcode: request.Body.Shortcode,
			Name:      request.Body.Name,
			Flags:     request.Body.Flags,
			Idcountry: request.Body.Idcountry,
			Depotno:   mapToNullString(request.Body.Depotno),
		}

		dbNewCustodian, err := queriesWithTx.InsertCustodian(ctx, insertParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		if len(request.Body.Exchanges) > 0 {

			//-----------------------------------------------------------------
			// Step A: Find default-exchange in collection
			//-----------------------------------------------------------------

			var index int = -1

			if request.Body.Idexchangedefault != nil {

				index = slices.IndexFunc(request.Body.Exchanges,

					func(entry Exchange2CustodianListItem) bool {

						return entry.Idexchange == *request.Body.Idexchangedefault
					})
			}

			if index < 0 {

				rs.logger.ErrorContext(ctx, OperationId, "error", fmt.Errorf("no default-exchange found in array: %s", request.Body.Idexchangedefault))

				return nil, huma.Error500InternalServerError("Internal server error")
			}

			//-----------------------------------------------------------------
			// Step B: Insert default-exchange at position 0 into db
			//-----------------------------------------------------------------

			err = queriesWithTx.InsertExchange2Custodian(ctx, rbsdb.InsertExchange2CustodianParams{

				Idcustodian: dbNewCustodian.ID,
				Idexchange:  request.Body.Exchanges[index].Idexchange,
				Sequenceno:  0,
				Value1:      request.Body.Exchanges[index].Value1,
				Value2:      request.Body.Exchanges[index].Value2})

			if err != nil {

				rs.logger.ErrorContext(ctx, OperationId, "error", err)

				return nil, huma.Error500InternalServerError("Internal server error")
			}

			//-----------------------------------------------------------------
			// Step C: Remove default-exchange from collection
			//-----------------------------------------------------------------

			request.Body.Exchanges = slices.Delete(request.Body.Exchanges, index, index+1)

			//-----------------------------------------------------------------
			// Step D: Insert all remaining echanges to db
			//-----------------------------------------------------------------

			for sequenceno, exchange := range request.Body.Exchanges {

				err = queriesWithTx.InsertExchange2Custodian(ctx, rbsdb.InsertExchange2CustodianParams{

					Idcustodian: dbNewCustodian.ID,
					Idexchange:  exchange.Idexchange,
					Sequenceno:  int32(sequenceno + 1),
					Value1:      exchange.Value1,
					Value2:      exchange.Value2})

				if err != nil {

					rs.logger.ErrorContext(ctx, OperationId, "error", err)

					return nil, huma.Error500InternalServerError("Internal server error")
				}
			}
		}

		//---------------------------------------------------------------------

		tx.Commit(ctx)

		var idexchangedefault uuid.UUID = uuid.Nil

		if request.Body.Idexchangedefault != nil {

			idexchangedefault = *request.Body.Idexchangedefault
		}

		return &CustodianResponseCreate{
			Header: LocationHeader{Location: fmt.Sprintf("/custodian/id/%s", dbNewCustodian.ID.String())},
			Body:   mapDB2APICustodian(dbNewCustodian, request.Body.Exchange2CustodianList, idexchangedefault),
		}, nil

	}, describeEndpoint("createCustodian", "Create a new custodian"), defaultStatus(http.StatusCreated))

	//-------------------------------------------------------------------------

	huma.Delete(group, "/{id}", func(ctx context.Context, request *CustodianRequestId) (*struct{}, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		dbResult, err := rs.dbQueries.DeleteCustodian(ctx, request.ID)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := dbResult.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil

	}, describeEndpoint("deleteCustodian", "Delete a single custodian based on the id supplied")) // huma-Defaultstatus = http.StatusNoContent

	//-------------------------------------------------------------------------

	huma.Put(group, "/{id}", func(ctx context.Context, request *CustodianRequestUpdate) (*CustodianResponse, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		//---------------------------------------------------------------------
		// Initialize transaction
		//---------------------------------------------------------------------

		tx, err := rs.dbPool.Begin(ctx)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		defer tx.Rollback(ctx)

		queriesWithTx := rs.dbQueries.WithTx(tx)

		//---------------------------------------------------------------------

		updateParams := rbsdb.UpdateCustodianParams{

			ID:        request.ID,
			Shortcode: request.Body.Shortcode,
			Name:      request.Body.Name,
			Flags:     request.Body.Flags,
			Idcountry: request.Body.Idcountry,
			Depotno:   mapToNullString(request.Body.Depotno),
		}

		dbUpdatedCustodian, err := queriesWithTx.UpdateCustodian(ctx, updateParams)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("")
			}

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		//---------------------------------------------------------------------
		// Step A: Delete all records for given 'Idcustodian'
		//---------------------------------------------------------------------

		_, err = queriesWithTx.DeleteExchange2CustodianByIdcustodian(ctx, dbUpdatedCustodian.ID)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, huma.Error500InternalServerError("Internal server error")
		}

		if len(request.Body.Exchanges) > 0 {

			//-----------------------------------------------------------------
			// Step B: Find default-exchange in collection
			//-----------------------------------------------------------------

			var index int = -1

			if request.Body.Idexchangedefault != nil {

				index = slices.IndexFunc(request.Body.Exchanges,

					func(entry Exchange2CustodianListItem) bool {

						return entry.Idexchange == *request.Body.Idexchangedefault
					})
			}

			if index < 0 {

				rs.logger.ErrorContext(ctx, OperationId, "error", fmt.Errorf("no default-exchange found in array: %s", request.Body.Idexchangedefault))

				return nil, huma.Error500InternalServerError("Internal server error")
			}

			//-----------------------------------------------------------------
			// Step C: Insert default-exchange at position 0 into db
			//-----------------------------------------------------------------

			err = queriesWithTx.InsertExchange2Custodian(ctx, rbsdb.InsertExchange2CustodianParams{

				Idcustodian: dbUpdatedCustodian.ID,
				Idexchange:  request.Body.Exchanges[index].Idexchange,
				Sequenceno:  0,
				Value1:      request.Body.Exchanges[index].Value1,
				Value2:      request.Body.Exchanges[index].Value2})

			if err != nil {

				rs.logger.ErrorContext(ctx, OperationId, "error", err)

				return nil, huma.Error500InternalServerError("Internal server error")
			}

			//-----------------------------------------------------------------
			// Step D: Remove default-exchange from collection
			//-----------------------------------------------------------------

			request.Body.Exchanges = slices.Delete(request.Body.Exchanges, index, index+1)

			//-----------------------------------------------------------------
			// Step E: Insert all remaining echanges to db
			//-----------------------------------------------------------------

			for sequenceno, exchange := range request.Body.Exchanges {

				err = queriesWithTx.InsertExchange2Custodian(ctx, rbsdb.InsertExchange2CustodianParams{

					Idcustodian: dbUpdatedCustodian.ID,
					Idexchange:  exchange.Idexchange,
					Sequenceno:  int32(sequenceno + 1),
					Value1:      exchange.Value1,
					Value2:      exchange.Value2})

				if err != nil {

					rs.logger.ErrorContext(ctx, OperationId, "error", err)

					return nil, huma.Error500InternalServerError("Internal server error")
				}
			}
		}

		//---------------------------------------------------------------------

		tx.Commit(ctx)

		var idexchangedefault uuid.UUID = uuid.Nil

		if request.Body.Idexchangedefault != nil {

			idexchangedefault = *request.Body.Idexchangedefault
		}

		return &CustodianResponse{Body: mapDB2APICustodian(dbUpdatedCustodian, request.Body.Exchange2CustodianList, idexchangedefault)}, nil

	}, describeEndpoint("updateCustodian", "Update an existing custodian based on the id supplied")) // huma-Defaultstatus = http.StatusOk
}

//-----------------------------------------------------------------------------

func mapDB2APICustodianListItem(record rbsdb.Custodian) CustodianListItem {

	return CustodianListItem{
		CustodianId:        CustodianId{ID: record.ID},
		CustodianShortcode: CustodianShortcode{Shortcode: record.Shortcode},
		CustodianName:      CustodianName{Name: record.Name},
	}
}

func mapDB2APICustodian(record rbsdb.Custodian, exchange2custodianlist Exchange2CustodianList, idexchangedefault uuid.UUID) Custodian {

	var custodian = &Custodian{
		CustodianId:            CustodianId{ID: record.ID},
		CustodianShortcode:     CustodianShortcode{Shortcode: record.Shortcode},
		CustodianName:          CustodianName{Name: record.Name},
		CustodianFlags:         CustodianFlags{Flags: record.Flags},
		CustodianIdCountry:     CustodianIdCountry{Idcountry: record.Idcountry},
		Exchange2CustodianList: Exchange2CustodianList{Exchanges: exchange2custodianlist.Exchanges},
	}

	if idexchangedefault != uuid.Nil {

		custodian.Exchange2CustodianIdExchangeDefault.Idexchangedefault = &idexchangedefault
	}

	return *custodian
}
