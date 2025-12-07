# Test script to demonstrate introspection features

def test_introspection():
    print("=== Testing Introspection Features ===\n")

    # Test dir() - list all available modules
    print("Available modules:")
    modules = dir()
    for mod in modules:
        print("  - " + mod)

    print("\n" + "="*50 + "\n")

    # Test help() - get help on all modules
    print(help())

    print("="*50 + "\n")

    # Test help(module) - get help on satellite module
    print(help(satellite))

    print("="*50 + "\n")

    # Test dir(module) - list functions in satellite module
    print("Functions in satellite module:")
    funcs = dir(satellite)
    for func in funcs:
        print("  - " + func)

# Run the test
test_introspection()
