/*
 * Data Catalog Service - Asset Details
 *
 * API version: 1.0.0
 * Based on code Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

// CHANGE-FROM-GENERATED-CODE: All code in this file is different from auto-generated code.
// This code is specific for working with Apache OpenMetadata

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	client "github.com/fybrik/datacatalog-go-client"
	models "github.com/fybrik/datacatalog-go-models"
	api "github.com/fybrik/datacatalog-go/go"
)

type ApacheApiService struct {
	hostname string
	port     string
	username string
	password string
}

// NewApacheApiService creates a new api service
func NewApacheApiService(conf map[interface{}]interface{}) OpenMetadataApiServicer {
	return &ApacheApiService{conf["openmetadata_hostname"].(string),
		strconv.Itoa(conf["openmetadata_port"].(int)),
		conf["openmetadata_username"].(string),
		conf["openmetadata_password"].(string)}
}

func waitUntilAssetIsDiscovered(ctx context.Context, c *client.APIClient, name string) bool {
	count := 0
	for {
		fmt.Println("running GetByName5")
		_, _, err := c.TablesApi.GetByName5(ctx, name).Execute()
		if err == nil {
			fmt.Println("Found the table!")
			return true
		} else {
			fmt.Println("Could not find the table. Let's try again")
		}

		if count == 10 {
			break
		}
		count++
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Too many retries. Could not find table. Giving up")
	return false
}

// CreateAsset - This REST API writes data asset information to the data catalog configured in fybrik
func (s *ApacheApiService) CreateAsset(ctx context.Context,
	xRequestDatacatalogWriteCred string,
	createAssetRequest models.CreateAssetRequest) (api.ImplResponse, error) {

	if createAssetRequest.Details.Connection.Name != "mysql" {
		return api.Response(400, nil), errors.New("currently, we only support the mysql connection")
	}

	conf := client.NewConfiguration()
	c := client.NewAPIClient(conf)

	// Let us begin with checking whether the database service already exists
	// XXXXXXXXX

	// If does not exist, let us create database service
	connection := client.NewDatabaseConnection()
	connection.SetConfig(createAssetRequest.Details.GetConnection().AdditionalProperties["mysql"].(map[string]interface{}))
	createDatabaseService := client.NewCreateDatabaseService(*connection, createAssetRequest.DestinationCatalogID+"-mysql", "Mysql")

	createDatabaseService, r, err := c.ServicesApi.Create16(ctx).CreateDatabaseService(*createDatabaseService).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicesApi.Create16``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return api.Response(r.StatusCode, nil), err
	}

	// get database service ID by name
	databaseService, r, err := c.ServicesApi.GetByName15(ctx, createDatabaseService.Name).Execute()

	// Next let us create an ingestion pipeline
	sourceConfig := *client.NewSourceConfig()
	sourceConfig.SetConfig(map[string]interface{}{"type": "DatabaseMetadata"})
	newCreateIngestionPipeline := *client.NewCreateIngestionPipeline(*&client.AirflowConfig{},
		"pipeline-"+*createAssetRequest.DestinationAssetID,
		"metadata", *client.NewEntityReference(databaseService.Id, "databaseService"),
		sourceConfig)

	createIngestionPipeline, r, err := c.IngestionPipelinesApi.Create17(ctx).CreateIngestionPipeline(newCreateIngestionPipeline).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicesApi.Create16``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return api.Response(r.StatusCode, nil), err
	}

	ingestionFQN := *createIngestionPipeline.Service.FullyQualifiedName + ".\"" + createIngestionPipeline.Name + "\""
	// get ingestion pipeline ID by name
	ingestionPipeline, r, err := c.IngestionPipelinesApi.GetByName16(ctx, ingestionFQN).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IngestionPipelinesApi.GetByName16``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return api.Response(r.StatusCode, nil), err
	}

	// Let us deploy the ingestion pipeline
	ingestionPipeline, r, err = c.IngestionPipelinesApi.DeployIngestion(ctx, *ingestionPipeline.Id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IngestionPipelinesApi.DeployIngestion``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return api.Response(r.StatusCode, nil), err
	}
	ingestionPipeline, r, err = c.IngestionPipelinesApi.TriggerIngestion(ctx, *ingestionPipeline.Id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IngestionPipelinesApi.TriggerIngestion``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return api.Response(r.StatusCode, nil), err
	}

	assetID := *createIngestionPipeline.Service.FullyQualifiedName + "." + *createAssetRequest.DestinationAssetID
	success := waitUntilAssetIsDiscovered(ctx, c, assetID)

	if success {
		return api.Response(201, api.CreateAssetResponse{}), nil
	} else {
		return api.Response(http.StatusNotImplemented, nil), errors.New("Could not find table " + assetID)
	}
}

// DeleteAsset - This REST API deletes data asset
func (s *ApacheApiService) DeleteAsset(ctx context.Context, xRequestDatacatalogCred string, deleteAssetRequest api.DeleteAssetRequest) (api.ImplResponse, error) {
	// TODO - update DeleteAsset with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, DeleteAssetResponse{}) or use other options such as http.Ok ...
	//return Response(200, DeleteAssetResponse{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return Response(401, nil),nil

	return api.Response(http.StatusNotImplemented, nil), errors.New("DeleteAsset method not implemented")
}

// GetAssetInfo - This REST API gets data asset information from the data catalog configured in fybrik for the data sets indicated in FybrikApplication yaml
func (s *ApacheApiService) GetAssetInfo(ctx context.Context, xRequestDatacatalogCred string, getAssetRequest api.GetAssetRequest) (api.ImplResponse, error) {
	conf := client.NewConfiguration()
	c := client.NewAPIClient(conf)

	assetID := getAssetRequest.AssetID

	//fields := "tableConstraints,tablePartition,usageSummary,owner,profileSample,customMetrics,tags,followers,joins,sampleData,viewDefinition,tableProfile,location,tableQueries,dataModel,tests" // string | Fields requested in the returned resource (optional)
	fields := "tags"
	include := "non-deleted" // string | Include all, deleted, or non-deleted entities. (optional) (default to "non-deleted")
	respAsset, r, err := c.TablesApi.GetByName5(ctx, assetID).Fields(fields).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TablesApi.GetByName5``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return api.Response(400, nil), err
	}

	serviceType := strings.ToLower(*respAsset.ServiceType)

	ret := &models.GetAssetResponse{}
	ret.Details.Connection.Name = serviceType
	dataFormat := "SQL"
	ret.Details.DataFormat = &dataFormat

	respService, r, err := c.ServicesApi.Get19(ctx, respAsset.Service.Id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicesApi.Get19``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return api.Response(400, nil), err
	}

	config := respService.Connection.GetConfig()

	additionalProperties := make(map[string]interface{})
	additionalProperties[serviceType] = config
	ret.Details.Connection.AdditionalProperties = additionalProperties
	ret.ResourceMetadata.Name = respAsset.FullyQualifiedName

	ret.Credentials = config["username"].(string) + ":" + config["password"].(string)

	for _, s := range respAsset.Columns {
		tags := make(map[string]interface{})
		for _, t := range s.Tags {
			tags[t.TagFQN] = "true"
		}
		ret.ResourceMetadata.Columns = append(ret.ResourceMetadata.Columns, models.ResourceColumn{Name: s.Name, Tags: tags})
	}

	tags := make(map[string]interface{})
	for _, s := range respAsset.Tags {
		tags[s.TagFQN] = "true"
	}
	ret.ResourceMetadata.Tags = tags

	return api.Response(200, ret), nil
}

// UpdateAsset - This REST API updates data asset information in the data catalog configured in fybrik
func (s *ApacheApiService) UpdateAsset(ctx context.Context, xRequestDatacatalogUpdateCred string, updateAssetRequest api.UpdateAssetRequest) (api.ImplResponse, error) {
	// TODO - update UpdateAsset with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, UpdateAssetResponse{}) or use other options such as http.Ok ...
	//return Response(200, UpdateAssetResponse{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return Response(401, nil),nil

	return api.Response(http.StatusNotImplemented, nil), errors.New("UpdateAsset method not implemented")
}
