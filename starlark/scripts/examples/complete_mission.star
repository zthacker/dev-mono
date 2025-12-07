# Complete mission example demonstrating all modules

def complete_mission():
    print("=== Complete Mission Example ===\n")

    # Step 1: Pre-flight checks using satellite telemetry
    print("Step 1: Pre-flight checks")
    battery = satellite.getTLM("battery_level")
    temp = satellite.getTLM("temperature")

    print("  Battery: " + str(battery) + "%")
    print("  Temperature: " + str(temp) + "°C")

    if battery < 20:
        print("  ABORT: Insufficient battery")
        return "ABORTED"

    print("  Pre-flight checks: PASSED\n")

    # Step 2: Ground station tracking
    print("Step 2: Ground station tracking")
    track_success = ground.track("SAT-001", 300)
    if not track_success:
        print("  Tracking failed")
        return "FAILED"
    print("  Tracking: ACTIVE\n")

    # Step 3: Send commands to satellite
    print("Step 3: Sending commands")
    satellite.sendCMD("POWER_ON", {"subsystem": "payload"})
    system.wait(2)
    satellite.sendCMD("SET_TRANSMITTER", {"power": "high", "freq": 2200})
    print("  Commands sent successfully\n")

    # Step 4: Data collection and processing
    print("Step 4: Data collection")
    raw_telemetry = "telemetry_data_stream"
    processed_data = data.process(raw_telemetry)
    is_valid = data.validate("telemetry_schema", processed_data)

    if is_valid:
        print("  Data collected and validated\n")
    else:
        print("  Data validation failed\n")
        return "DATA_ERROR"

    # Step 5: Schedule next pass
    print("Step 5: Scheduling next pass")
    next_pass_config = {
        "start_time": "2025-12-08T14:00:00Z",
        "duration": 600
    }
    schedule_id = ground.schedule("PASS-789", next_pass_config)
    print("  Next pass scheduled: " + schedule_id + "\n")

    print("=== Mission Complete: SUCCESS ===")
    return "SUCCESS"

# Execute the mission
result = complete_mission()
print("\nFinal Status: " + result)
