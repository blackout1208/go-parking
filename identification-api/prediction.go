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

func (p *Prediction) getPlatesIMG() []bbox {
	var platesIMG []bbox

	for i, confidence := range p.Confidences {
		if confidence < 0.2 || p.DisplayNames[i] != _licensePlateLabel {
			continue
		}

		platesIMG = append(platesIMG, bbox{
			xmin: p.Bboxes[i][0],
			xmax: p.Bboxes[i][1],
			ymin: p.Bboxes[i][2],
			ymax: p.Bboxes[i][3],
		})
	}

	return platesIMG
}
