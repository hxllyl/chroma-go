package chroma_go

import (
	"context"
	"encoding/json"
	"fmt"
	openapiclient "github.com/amikos-tech/chroma-go/swagger"
	"log"
	"reflect"
	"strings"
)

type ClientConfiguration struct {
	BasePath          string            `json:"basePath,omitempty"`
	DefaultHeaders    map[string]string `json:"defaultHeader,omitempty"`
	EmbeddingFunction EmbeddingFunction `json:"embeddingFunction,omitempty"`
}

type EmbeddingFunction interface {
	CreateEmbedding(documents []string) ([][]float32, error)
	CreateEmbeddingWithModel(documents []string, model string) ([][]float32, error)
}

func MapToApi(inmap map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range inmap {
		result[k] = v
	}
	return result
}

func MapListToApi(inmap []map[string]string) []map[string]interface{} {
	result := make([]map[string]interface{}, len(inmap))
	for i, v := range inmap {
		result[i] = MapToApi(v)
	}
	return result
}

func MapFromApi(inmap map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range inmap {
		result[k] = v.(string)
	}
	return result
}

func MapListFromApi(inmap []map[string]interface{}) []map[string]string {
	result := make([]map[string]string, len(inmap))
	for i, v := range inmap {
		result[i] = MapFromApi(v)
	}
	return result
}

// Client represents the ChromaDB Client
type Client struct {
	ApiClient *openapiclient.APIClient
}

func NewClient(basePath string) *Client {
	configuration := openapiclient.NewConfiguration()
	_, err := configuration.Servers.URL(0, map[string]string{"basePath": basePath})
	if err != nil {
		return nil
	}
	apiClient := openapiclient.NewAPIClient(configuration)
	return &Client{
		ApiClient: apiClient,
	}
}

func (c *Client) GetCollection(collectionName string, embeddingFunction EmbeddingFunction) (*Collection, error) {
	col, httpResp, err := c.ApiClient.DefaultApi.GetCollection(context.Background(), collectionName).Execute()
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("error getting collection: %v", httpResp)
	}
	return NewCollection(c.ApiClient, col.Id, col.Name, col.Metadata, embeddingFunction), nil
}

func (c *Client) Heartbeat() (map[string]int32, error) {
	resp, httpResp, err := c.ApiClient.DefaultApi.Heartbeat(context.Background()).Execute()
	fmt.Printf("Heartbeat: %v\n", httpResp)
	return resp, err
}

type DistanceFunction string

const (
	L2     DistanceFunction = "l2"
	COSINE DistanceFunction = "cosine"
	IP     DistanceFunction = "ip"
)

func GetStringTypeOfEmbeddingFunction(ef EmbeddingFunction) string {
	typ := reflect.TypeOf(ef)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem() // Dereference if it's a pointer
	}
	return typ.String()
}

func (c *Client) CreateCollection(collectionName string, metadata map[string]interface{}, createOrGet bool, embeddingFunction EmbeddingFunction, distanceFunction DistanceFunction) (*Collection, error) {
	_metadata := metadata

	if _metadata == nil || len(_metadata) == 0 {
		_metadata = make(map[string]interface{})
	}
	if _metadata["embedding_function"] == "" {
		_metadata["embedding_function"] = GetStringTypeOfEmbeddingFunction(embeddingFunction)
	}
	if distanceFunction == "" {
		_metadata["hnsw:space"] = strings.ToLower(string(L2))
	} else {
		_metadata["hnsw:space"] = strings.ToLower(string(distanceFunction))
	}

	col := openapiclient.CreateCollection{
		Name:        collectionName,
		GetOrCreate: &createOrGet,
		Metadata:    _metadata,
	}
	resp, httpResp, err := c.ApiClient.DefaultApi.CreateCollection(context.Background()).CreateCollection(col).Execute()
	if err != nil {
		return nil, err
	}
	fmt.Printf("CreateCollection: %v\n", httpResp.Body)
	respJSON, _ := json.Marshal(resp)
	fmt.Println(string(respJSON))
	mtd := resp.Metadata
	return NewCollection(c.ApiClient, resp.Id, resp.Name, mtd, embeddingFunction), nil
}

func (c *Client) DeleteCollection(collectionName string) (*Collection, error) {
	_, httpResp, gcerr := c.ApiClient.DefaultApi.GetCollection(context.Background(), collectionName).Execute()
	if gcerr != nil {
		log.Fatal(httpResp, gcerr)
		return nil, gcerr
	}
	deletedCol, httpResp, err := c.ApiClient.DefaultApi.DeleteCollection(context.Background(), collectionName).Execute()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return NewCollection(c.ApiClient, deletedCol.Id, deletedCol.Name, deletedCol.Metadata, nil), nil
}

func (c *Client) Reset() (bool, error) {
	resp, httpResp, err := c.ApiClient.DefaultApi.Reset(context.Background()).Execute()
	fmt.Printf("Reset: %v\n", httpResp)
	return resp, err
}

func (c *Client) ListCollections() ([]*Collection, error) {
	req := c.ApiClient.DefaultApi.ListCollections(context.Background())
	resp, httpResp, err := req.Execute()
	fmt.Printf("ListCollections: %v\n", httpResp)
	if err != nil {
		return nil, err
	}
	collections := make([]*Collection, len(resp))
	for i, col := range resp {
		collections[i] = NewCollection(c.ApiClient, col.Id, col.Name, col.Metadata, nil)
	}
	return collections, nil
}

func (c *Client) Version() (string, error) {
	req := c.ApiClient.DefaultApi.Version(context.Background())
	resp, httpResp, err := req.Execute()
	fmt.Printf("Version: %v\n", httpResp)
	return resp, err
}

type CollectionData struct {
	Ids       []string                 `json:"ids,omitempty"`
	Documents []string                 `json:"documents,omitempty"`
	Metadatas []map[string]interface{} `json:"metadatas,omitempty"`
}

type Collection struct {
	Name              string
	EmbeddingFunction EmbeddingFunction
	ApiClient         *openapiclient.APIClient
	Metadata          map[string]interface{}
	id                string
	CollectionData    *CollectionData
}

func NewCollection(apiClient *openapiclient.APIClient, id string, name string, metadata map[string]interface{}, embeddingFunction EmbeddingFunction) *Collection {
	return &Collection{
		Name:              name,
		EmbeddingFunction: embeddingFunction,
		ApiClient:         apiClient,
		Metadata:          metadata,
		id:                id,
	}
}

func (c *Collection) Add(embeddings [][]float32, metadatas []map[string]interface{}, documents []string, ids []string) (*Collection, error) {
	req := c.ApiClient.DefaultApi.Add(context.Background(), c.id)

	var _embeddings []interface{}

	if len(embeddings) == 0 {
		embds, embErr := c.EmbeddingFunction.CreateEmbedding(documents)
		if embErr != nil {
			return c, embErr
		}
		_embeddings = ConvertEmbeds(embds)
	} else {
		_embeddings = ConvertEmbeds(embeddings)
	}

	req.AddEmbedding(openapiclient.AddEmbedding{
		Embeddings: _embeddings,
		Metadatas:  metadatas,
		Documents:  documents,
		Ids:        ids,
	})

	_, httpResp, err := req.Execute()

	if err != nil {
		return c, err
	}
	fmt.Printf("Add: %v\n", httpResp)
	return c, nil
}

func (c *Collection) Upsert(embeddings [][]float32, metadatas []map[string]interface{}, documents []string, ids []string) (*Collection, error) {
	req := c.ApiClient.DefaultApi.Upsert(context.Background(), c.id)

	var _embeddings []interface{}

	if len(embeddings) == 0 {
		embds, embErr := c.EmbeddingFunction.CreateEmbedding(documents)
		if embErr != nil {
			return c, embErr
		}
		_embeddings = ConvertEmbeds(embds)
	} else {
		_embeddings = ConvertEmbeds(embeddings)
	}

	req.AddEmbedding(openapiclient.AddEmbedding{
		Embeddings: _embeddings,
		Metadatas:  metadatas,
		Documents:  documents,
		Ids:        ids,
	})

	_, httpResp, err := req.Execute()

	if err != nil {
		return c, err
	}
	fmt.Printf("Add: %v\n", httpResp)
	return c, nil
}

func (c *Collection) Get(where map[string]interface{}, whereDocuments map[string]interface{}, ids []string) (*Collection, error) {
	req := c.ApiClient.DefaultApi.Get(context.Background(), c.id)
	req.GetEmbedding(openapiclient.GetEmbedding{
		Ids:           ids,
		Where:         where,
		WhereDocument: whereDocuments,
	})
	//{
	//	Ids:           ids,
	//	Where:         where,
	//	WhereDocument: whereDocuments,
	//}

	cd, httpResp, err := req.Execute()

	if err != nil {
		return c, err
	}
	cdata := CollectionData{
		Ids:       cd.Ids,
		Documents: cd.Documents,
		Metadatas: getMetadatasListFromAPI(cd.Metadatas),
	}
	c.CollectionData = &cdata
	fmt.Printf("Add: %v\n", httpResp)
	return c, nil
}

type QueryEnum string

const (
	documents  QueryEnum = "documents"
	embeddings QueryEnum = "embeddings"
	metadatas  QueryEnum = "metadatas"
	distances  QueryEnum = "distances"
)

type QueryResults struct {
	Documents [][]string                 `json:"documents,omitempty"`
	Ids       [][]string                 `json:"ids,omitempty"`
	Metadatas [][]map[string]interface{} `json:"metadatas,omitempty"`
	Distances [][]float32                `json:"distances,omitempty"`
}

func getMetadatasListFromAPI(metadatas []map[string]openapiclient.MetadatasInnerValue) []map[string]interface{} {
	// Initialize the result slice
	result := make([]map[string]interface{}, len(metadatas))
	// Iterate over the inner map
	for j, metadataMap := range metadatas {
		resultMap := make(map[string]interface{})
		for key, value := range metadataMap {
			// Convert MetadatasInnerValue to interface{}
			var rawValue interface{}
			b, e := value.MarshalJSON()
			if e != nil {
				rawValue = nil
			}
			rawValue = b
			// Store in the result map
			resultMap[key] = rawValue
		}
		result[j] = resultMap
	}

	return result
}

func getMetadatasFromAPI(metadatas [][]map[string]openapiclient.MetadatasInnerValue) [][]map[string]interface{} {
	// Initialize the result slice
	result := make([][]map[string]interface{}, len(metadatas))

	// Iterate over the outer slice
	for i, outerItem := range metadatas {
		result[i] = make([]map[string]interface{}, len(outerItem))

		// Iterate over the inner map
		for j, metadataMap := range outerItem {
			resultMap := make(map[string]interface{})
			for key, value := range metadataMap {
				// Convert MetadatasInnerValue to interface{}
				var rawValue interface{}
				b, e := value.MarshalJSON()
				if e != nil {
					rawValue = nil
				}
				rawValue = b
				// Store in the result map
				resultMap[key] = rawValue
			}
			result[i][j] = resultMap
		}
	}

	return result
}
func ConvertEmbeds(embeds [][]float32) []interface{} {
	_embeddings := make([]interface{}, len(embeds))
	for i, v := range embeds {
		_embeddings[i] = v
	}
	return _embeddings
}
func (c *Collection) Query(queryTexts []string, nResults int32, where map[string]interface{}, whereDocuments map[string]interface{}, include []QueryEnum) (*QueryResults, error) {
	_includes := make([]openapiclient.IncludeInner, len(include))
	for i, v := range include {
		var inr = openapiclient.IncludeInner{}
		inr.UnmarshalJSON([]byte(v))
		_includes[i] = inr
	}

	embds, embErr := c.EmbeddingFunction.CreateEmbedding(queryTexts)
	if embErr != nil {
		return nil, embErr
	}

	nreq := c.ApiClient.DefaultApi.GetNearestNeighbors(context.Background(), c.id)
	nreq.QueryEmbedding(openapiclient.QueryEmbedding{
		Where:           where,
		WhereDocument:   whereDocuments,
		NResults:        &nResults,
		Include:         _includes,
		QueryEmbeddings: ConvertEmbeds(embds),
	})
	qr, httpResp, err := nreq.Execute()
	if err != nil {
		return nil, err
	}
	qresults := QueryResults{
		Documents: qr.Documents,
		Ids:       qr.Ids,
		Metadatas: getMetadatasFromAPI(qr.Metadatas),
		Distances: qr.Distances,
	}
	fmt.Printf("Add: %v\n", httpResp)
	return &qresults, nil

}

func (c *Collection) Count() (int32, error) {
	req := c.ApiClient.DefaultApi.Count(context.Background(), c.id)

	cd, httpResp, err := req.Execute()

	if err != nil {
		return -1, err
	}

	fmt.Printf("Count: %v\n", httpResp)

	return cd, nil
}

func (c *Collection) Update(newName string, newMetadata map[string]interface{}) (*Collection, error) {
	req := c.ApiClient.DefaultApi.UpdateCollection(context.Background(), c.id)
	req.UpdateCollection(openapiclient.UpdateCollection{
		NewName:     &newName,
		NewMetadata: newMetadata,
	})

	col, httpResp, err := req.Execute()
	if err != nil {
		log.Fatal(httpResp, err)
		return c, err
	}
	c.Name = col.Name
	c.Metadata = col.Metadata
	return c, nil
}

func (c *Collection) Delete(ids []string, where map[string]interface{}, whereDocuments map[string]interface{}) ([]string, error) {
	req := c.ApiClient.DefaultApi.Delete(context.Background(), c.id)
	req.DeleteEmbedding(openapiclient.DeleteEmbedding{
		Where:         where,
		WhereDocument: whereDocuments,
		Ids:           ids,
	})

	dr, httpResp, err := req.Execute()
	if err != nil {
		log.Fatal(httpResp, err)
		return nil, err
	}
	return dr, nil

}
