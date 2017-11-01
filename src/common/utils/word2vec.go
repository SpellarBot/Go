package utils

import "fmt"

/*
Options:
Parameters for training:
	-train <file>
		Use text data from <file> to train the model
	-output <file>
		Use <file> to save the resulting word vectors / word clusters
	-size <int>
		Set size of word vectors; default is 100
	-window <int>
		Set max skip length between words; default is 5
	-sample <float>
		Set threshold for occurrence of words. Those that appear with higher frequency in the training data will be randomly down-sampled; default is 0 (off), useful value is 1e-5
	-hs <int>
		Use Hierarchical Softmax; default is 1 (0 = not used)
	-negative <int>
		Number of negative examples; default is 0, common values are 5 - 10 (0 = not used)
	-threads <int>
		Use <int> threads (default 1)
	-min-count <int>
		This will discard words that appear less than <int> times; default is 5
	-alpha <float>
		Set the starting learning rate; default is 0.025
	-classes <int>
		Output word classes rather than word vectors; default number of classes is 0 (vectors are written)
	-debug <int>
		Set the debug mode (default = 2 = more info during training)
	-binary <int>
		Save the resulting vectors in binary moded; default is 0 (off)
	-save-vocab <file>
		The vocabulary will be saved to <file>
	-read-vocab <file>
		The vocabulary will be read from <file>, not constructed from the training data
	-cbow <int>
		Use the continuous bag of words model; default is 0 (skip-gram model)
*/


type Word2vecRequest struct {
	Bin            string
	TrainFileName  string
	OutputFileName string
	HS             bool
	Negative       int
	CBOW           bool
	DimSize        int
	Window         int
	Sample         string
	Threads        int
	MinCount       int
	Alpha          string
	Classes        int
	Debug          int
	Binary         bool
	SaveVocab      string
	ReadVocab      string
}

func Word2vec(r Word2vecRequest) error {
	if (r.Bin == "") || (r.TrainFileName == "") || (r.OutputFileName == "") {
		return fmt.Errorf("[Word2vec] Error! Must Set TrainFileName And OutputFileName!")
	}

	params := []string{}

	params = append(params, "-train")
	params = append(params, r.TrainFileName)

	params = append(params, "-output")
	params = append(params, r.OutputFileName)

	// Use Hierarchical Softmax or NOT
	params = append(params, "-hs")
	if r.HS {
		params = append(params, "1")
	} else {
		params = append(params, "0")
	}

	// User Negative Sampling or NOT
	if r.Negative > 0 {
		params = append(params, "-negative")
		params = append(params, fmt.Sprintf("%d", r.Negative))
	}

	// Use CBOW or Skip-Gram
	params = append(params, "-cbow")
	if r.CBOW {
		params = append(params, "1")
	} else {
		params = append(params, "0")
	}

	if r.DimSize > 0 {
		params = append(params, "-size")
		params = append(params, fmt.Sprintf("%d", r.DimSize))
	}

	if r.Window > 0 {
		params = append(params, "-window")
		params = append(params, fmt.Sprintf("%d", r.Window))
	}

	if r.Sample != "" {
		params = append(params, "-sample")
		params = append(params, r.Sample)
	}

	if r.Threads <= 0 {
		r.Threads = 1
	}
	params = append(params, "-threads")
	params = append(params, fmt.Sprintf("%d", r.Threads))

	if r.MinCount > 0 {
		params = append(params, "-min-count")
		params = append(params, fmt.Sprintf("%d", r.MinCount))
	}

	if r.Alpha != "" {
		params = append(params, "-alpha")
		params = append(params, r.Alpha)
	}

	if r.Classes > 0 {
		params = append(params, "-classes")
		params = append(params, fmt.Sprintf("%d", r.Classes))
	}

	if r.Debug > 0 {
		params = append(params, "-debug")
		params = append(params, fmt.Sprintf("%d", r.Debug))
	}

	params = append(params, "-binary")
	if r.Binary {
		params = append(params, "1")
	} else {
		params = append(params, "0")
	}

	if r.SaveVocab != "" {
		params = append(params, "-save-vocab")
		params = append(params, r.SaveVocab)
	}

	if r.ReadVocab != "" {
		params = append(params, "-read-vocab")
		params = append(params, r.ReadVocab)
	}

	_, err := Exec(r.Bin, params...)
	return err
}