package pos

import (
	"bytes"
	"encoding/gob"
)

/* Feature Gob interface */

func (sf singleFeature) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(sf.featureType); err != nil {
		return nil, err
	}

	if err := encoder.Encode(sf.value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (sf *singleFeature) GobDecode(buf []byte) error {
	b := bytes.NewBuffer(buf)

	decoder := gob.NewDecoder(b)

	if err := decoder.Decode(&sf.featureType); err != nil {
		return err
	}

	if err := decoder.Decode(&sf.value); err != nil {
		return err
	}

	return nil
}

func (tf tupleFeature) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(tf.featureType); err != nil {
		return nil, err
	}

	if err := encoder.Encode(tf.value1); err != nil {
		return nil, err
	}

	if err := encoder.Encode(tf.value2); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tf *tupleFeature) GobDecode(buf []byte) error {
	b := bytes.NewBuffer(buf)

	decoder := gob.NewDecoder(b)

	if err := decoder.Decode(&tf.featureType); err != nil {
		return err
	}

	if err := decoder.Decode(&tf.value1); err != nil {
		return err
	}

	if err := decoder.Decode(&tf.value2); err != nil {
		return err
	}

	return nil
}

/* fctuple Gob Interface */
func (fc fctuple) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(&fc.feature); err != nil {
		return nil, err
	}

	if err := encoder.Encode(fc.POSTag); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (fc *fctuple) GobDecode(buf []byte) error {
	b := bytes.NewBuffer(buf)

	decoder := gob.NewDecoder(b)
	if err := decoder.Decode(&fc.feature); err != nil {
		return err
	}

	if err := decoder.Decode(&fc.POSTag); err != nil {
		return err
	}
	return nil
}

/* Perceptron Gob Interface */

func (p *perceptron) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	// if err := encoder.Encode(&p.weights); err != nil {
	// 	return nil, err
	// }

	if err := encoder.Encode(&p.weightsSF); err != nil {
		return nil, err
	}
	if err := encoder.Encode(&p.weightsTF); err != nil {
		return nil, err
	}

	if err := encoder.Encode(&p.totals); err != nil {
		return nil, err
	}

	if err := encoder.Encode(&p.steps); err != nil {
		return nil, err
	}

	if err := encoder.Encode(p.instancesSeen); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *perceptron) GobDecode(buf []byte) error {
	b := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(b)

	// if err := decoder.Decode(&p.weights); err != nil {
	// 	return err
	// }

	if err := decoder.Decode(&p.weightsSF); err != nil {
		return err
	}

	if err := decoder.Decode(&p.weightsTF); err != nil {
		return err
	}

	if err := decoder.Decode(&p.totals); err != nil {
		return err
	}

	if err := decoder.Decode(&p.steps); err != nil {
		return err
	}

	if err := decoder.Decode(&p.instancesSeen); err != nil {
		return err
	}

	return nil
}

func init() {
	gob.Register(singleFeature{})
	gob.Register(tupleFeature{})
}
