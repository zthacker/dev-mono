# example starlark script for a mission

def run_pre_pass(target_name):
    print("Starting pre-pass for target: " + target_name)

    # check telem
    battery = satellite.getTLM("battery_level")
    if battery < 20:
        return "ABORT: Low Battery (" + str(battery) + "V)"    
    print("Battery level sufficient: " + str(battery) + "%")

    # send cmd
    # ops can send cmds using a dictionary for args
    satellite.sendCMD("SET_TRANSMITTER", {"power": "high", "freq": 2200})

    # wait for cmd to complete
    system.wait(2)

    return "SUCCESS: Transmitter set for mission."