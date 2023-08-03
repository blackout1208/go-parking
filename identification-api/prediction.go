package main

const (
	_licensePlateLabel = "license_plate"
)

type Prediction struct {
	// https://cloud.google.com/vertex-ai/docs/image-data/object-detection/interpret-results
	// "bboxes": [ [xMin, xMax, yMin, yMax], ...]
	Bboxes       [][]float64 `json:"bboxes"`
	Confidences  []float64   `json:"confidences"`
	DisplayNames []string    `json:"displayNames"`
}

func extractPrediction(respMap map[string]interface{}) Prediction {
	var prediction Prediction

	for _, item := range respMap["confidences"].([]interface{}) {
		prediction.Confidences = append(prediction.Confidences, item.(float64))
	}

	for _, item := range respMap["displayNames"].([]interface{}) {
		prediction.DisplayNames = append(prediction.DisplayNames, item.(string))
	}

	for _, item := range respMap["bboxes"].([]interface{}) {
		var bbox []float64
		for _, item2 := range item.([]interface{}) {
			bbox = append(bbox, item2.(float64))
		}
		prediction.Bboxes = append(prediction.Bboxes, bbox)
	}

	return prediction
}
