package estimation

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// ExtendedKalmanFilter implements the StateEstimator interface.
// This is where you'll implement the actual EKF algorithm.
type ExtendedKalmanFilter struct {
	// Process noise covariance
	Q *mat.Dense

	// Measurement noise covariance
	R *mat.Dense
}

// NewExtendedKalmanFilter creates a new EKF instance with default noise parameters.
func NewExtendedKalmanFilter() *ExtendedKalmanFilter {
	// Initialize with reasonable defaults (you'll tune these)
	Q := mat.NewDense(6, 6, nil)
	for i := 0; i < 6; i++ {
		Q.Set(i, i, 1e-6) // Small process noise
	}

	R := mat.NewDense(3, 3, nil)
	for i := 0; i < 3; i++ {
		R.Set(i, i, 100.0) // GPS measurement noise (meters^2)
	}

	return &ExtendedKalmanFilter{
		Q: Q,
		R: R,
	}
}

// Predict propagates the state forward in time using the dynamics model.
// For orbital mechanics, this typically uses Keplerian or SGP4 propagation.
func (ekf *ExtendedKalmanFilter) Predict(state *EKFState, dt float64) error {
	// TODO: Implement orbital dynamics propagation
	// For now, this is a simple constant-velocity model (placeholder)

	// State transition matrix F for constant velocity
	// x_k+1 = F * x_k
	F := mat.NewDense(6, 6, []float64{
		1, 0, 0, dt, 0, 0,
		0, 1, 0, 0, dt, 0,
		0, 0, 1, 0, 0, dt,
		0, 0, 0, 1, 0, 0,
		0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, 1,
	})

	// Predict state: x = F * x
	var newState mat.VecDense
	newState.MulVec(F, state.State)
	state.State = &newState

	// Predict covariance: P = F * P * F^T + Q
	var FP mat.Dense
	FP.Mul(F, state.Covariance)

	var FPFt mat.Dense
	FPFt.Mul(&FP, F.T())

	state.Covariance.Add(&FPFt, ekf.Q)

	return nil
}

// Update corrects the predicted state based on a GPS measurement.
func (ekf *ExtendedKalmanFilter) Update(state *EKFState, measurement *GPSMeasurement) error {
	// TODO: Implement full EKF update with potentially nonlinear measurement model
	// For GPS position measurements in ECI, the measurement model is linear: H = [I_3x3 | 0_3x3]

	// Measurement matrix H (measures position only, not velocity)
	H := mat.NewDense(3, 6, []float64{
		1, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0,
		0, 0, 1, 0, 0, 0,
	})

	// Innovation: y = z - H * x
	var Hx mat.VecDense
	Hx.MulVec(H, state.State)

	innovation := mat.NewVecDense(3, nil)
	innovation.SubVec(measurement.Position, &Hx)

	// Innovation covariance: S = H * P * H^T + R
	var HP mat.Dense
	HP.Mul(H, state.Covariance)

	var HPHt mat.Dense
	HPHt.Mul(&HP, H.T())

	S := mat.NewDense(3, 3, nil)
	S.Add(&HPHt, ekf.R)

	// Kalman gain: K = P * H^T * S^-1
	var PHt mat.Dense
	PHt.Mul(state.Covariance, H.T())

	var Sinv mat.Dense
	if err := Sinv.Inverse(S); err != nil {
		return err
	}

	K := mat.NewDense(6, 3, nil)
	K.Mul(&PHt, &Sinv)

	// Update state: x = x + K * innovation
	var Ky mat.VecDense
	Ky.MulVec(K, innovation)
	state.State.AddVec(state.State, &Ky)

	// Update covariance: P = (I - K * H) * P
	I := mat.NewDense(6, 6, nil)
	for i := 0; i < 6; i++ {
		I.Set(i, i, 1.0)
	}

	var KH mat.Dense
	KH.Mul(K, H)

	var IKH mat.Dense
	IKH.Sub(I, &KH)

	var newCov mat.Dense
	newCov.Mul(&IKH, state.Covariance)
	state.Covariance = &newCov

	return nil
}

// Name returns the algorithm name.
func (ekf *ExtendedKalmanFilter) Name() string {
	return "Extended Kalman Filter"
}

// SetProcessNoise allows tuning the process noise covariance.
func (ekf *ExtendedKalmanFilter) SetProcessNoise(Q *mat.Dense) {
	ekf.Q = Q
}

// SetMeasurementNoise allows tuning the measurement noise covariance.
func (ekf *ExtendedKalmanFilter) SetMeasurementNoise(R *mat.Dense) {
	ekf.R = R
}

// Helper function for orbital mechanics (placeholder for future implementation)
func computeOrbitalAcceleration(pos *mat.VecDense) *mat.VecDense {
	// TODO: Implement J2 perturbations, atmospheric drag, solar radiation pressure, etc.
	// For now, simple two-body problem: a = -μ * r / |r|^3

	const mu = 3.986004418e14 // Earth's gravitational parameter (m^3/s^2)

	x := pos.AtVec(0)
	y := pos.AtVec(1)
	z := pos.AtVec(2)

	r := math.Sqrt(x*x + y*y + z*z)
	r3 := r * r * r

	accel := mat.NewVecDense(3, []float64{
		-mu * x / r3,
		-mu * y / r3,
		-mu * z / r3,
	})

	return accel
}
