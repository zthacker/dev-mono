# Test script to verify satellite module is working

def test_satellite():
    print("=== Testing Satellite Module ===")

    # Get telemetry
    battery = satellite.getTLM("battery_level")
    print("Battery level: " + str(battery) + "%")

    # Check if battery is sufficient
    if battery < 20:
        print("ABORT: Low Battery (" + str(battery) + "V)")
    else:
        print("Battery level sufficient")

        # Send command
        satellite.sendCMD("SET_TRANSMITTER", {"power": "high", "freq": 2200})
        print("Transmitter command sent")

        # Wait for command to complete
        system.wait(1)

        print("SUCCESS: Operation complete")

# Run the test
test_satellite()
