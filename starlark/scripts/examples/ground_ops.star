# Ground station operations example

def ground_operations():
    print("=== Ground Station Operations ===")

    # Track a satellite
    sat_id = "SAT-001"
    duration = 300  # 5 minutes

    print("Initiating tracking for " + sat_id)
    success = ground.track(sat_id, duration)

    if success:
        print("Tracking successful")

        # Schedule next pass
        pass_config = {
            "start_time": "2025-12-08T10:00:00Z",
            "duration": 600,
            "elevation": "high"
        }

        schedule_id = ground.schedule("PASS-456", pass_config)
        print("Pass scheduled: " + schedule_id)
    else:
        print("Tracking failed")

# Run the example
ground_operations()
