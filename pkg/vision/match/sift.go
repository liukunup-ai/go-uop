package match

import (
	"gocv.io/x/gocv"
)

type siftMatcher struct {
	config *Config
}

func newSIFTMatcher(cfg *Config) *siftMatcher {
	return &siftMatcher{config: cfg}
}

func (m *siftMatcher) Name() string {
	return "sift"
}

func (m *siftMatcher) Find(screenshot, templateImg []byte) ([]*MatchResult, error) {
	if screenshot == nil || templateImg == nil {
		return []*MatchResult{}, nil
	}

	cfg := defaultConfig()
	if m.config != nil {
		cfg = m.config
	}

	img1, err := gocv.IMDecode(screenshot, gocv.IMReadGrayScale)
	if err != nil {
		return nil, err
	}
	defer img1.Close()

	img2, err := gocv.IMDecode(templateImg, gocv.IMReadGrayScale)
	if err != nil {
		return nil, err
	}
	defer img2.Close()

	sift := gocv.NewSIFT()
	defer sift.Close()

	kp1, des1 := sift.DetectAndCompute(img1, gocv.Mat{})
	defer des1.Close()

	kp2, des2 := sift.DetectAndCompute(img2, gocv.Mat{})
	defer des2.Close()

	if len(kp1) == 0 || len(kp2) == 0 {
		return []*MatchResult{}, nil
	}

	bf := gocv.NewBFMatcher()
	defer bf.Close()

	matches := bf.KnnMatch(des1, des2, 2)

	var goodMatches []gocv.DMatch
	for _, match := range matches {
		if len(match) > 1 {
			if match[0].Distance < 0.75*match[1].Distance {
				goodMatches = append(goodMatches, match[0])
			}
		}
	}

	if len(goodMatches) < 4 {
		return []*MatchResult{}, nil
	}

	srcPoints := make([]gocv.Point2f, len(goodMatches))
	dstPoints := make([]gocv.Point2f, len(goodMatches))
	for i, match := range goodMatches {
		srcPoints[i] = gocv.NewPoint2f(float32(kp1[match.QueryIdx].X), float32(kp1[match.QueryIdx].Y))
		dstPoints[i] = gocv.NewPoint2f(float32(kp2[match.TrainIdx].X), float32(kp2[match.TrainIdx].Y))
	}

	srcPointVec := gocv.NewPoint2fVectorFromPoints(srcPoints)
	defer srcPointVec.Close()
	dstPointVec := gocv.NewPoint2fVectorFromPoints(dstPoints)
	defer dstPointVec.Close()

	srcPointsMat := gocv.NewMatFromPoint2fVector(srcPointVec, true)
	defer srcPointsMat.Close()
	dstPointsMat := gocv.NewMatFromPoint2fVector(dstPointVec, true)
	defer dstPointsMat.Close()

	mask := gocv.NewMat()
	defer mask.Close()
	homography := gocv.FindHomography(srcPointsMat, dstPointsMat, gocv.HomographyMethodRANSAC, 3.0, &mask, 2000, 0.995)
	defer homography.Close()

	if homography.Empty() {
		return []*MatchResult{}, nil
	}

	homographyInv := homography.Clone()
	defer homographyInv.Close()
	homographyInv.Inv()

	threshold := cfg.Threshold
	if threshold <= 0 {
		threshold = 0.8
	}

	inliers := 0
	for i := 0; i < mask.Total(); i++ {
		if mask.GetIntAt(i, 0) == 1 {
			inliers++
		}
	}

	if float64(inliers)/float64(len(goodMatches)) < threshold {
		return []*MatchResult{}, nil
	}

	templateW := img2.Cols()
	templateH := img2.Rows()
	templateCorners := []gocv.Point2f{
		gocv.NewPoint2f(0, 0),
		gocv.NewPoint2f(float32(templateW), 0),
		gocv.NewPoint2f(float32(templateW), float32(templateH)),
		gocv.NewPoint2f(0, float32(templateH)),
	}

	srcVec := gocv.NewPoint2fVectorFromPoints(templateCorners)
	defer srcVec.Close()
	srcMat := gocv.NewMatFromPoint2fVector(srcVec, true)
	defer srcMat.Close()

	dstMat := gocv.NewMat()
	defer dstMat.Close()
	gocv.PerspectiveTransform(srcMat, &dstMat, homographyInv)

	minX := int(dstMat.GetFloatAt(0, 0))
	maxX := minX
	minY := int(dstMat.GetFloatAt(0, 1))
	maxY := minY
	for i := 0; i < 4; i++ {
		x := int(dstMat.GetFloatAt(i, 0))
		y := int(dstMat.GetFloatAt(i, 1))
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	score := float64(inliers) / float64(len(goodMatches))

	return []*MatchResult{
		{
			Found:  true,
			X:      minX,
			Y:      minY,
			Width:  maxX - minX,
			Height: maxY - minY,
			Score:  score,
		},
	}, nil
}

func (m *siftMatcher) DebugRender(screenshot []byte, results []*MatchResult) []byte {
	return debugRender(screenshot, results, m.config)
}
