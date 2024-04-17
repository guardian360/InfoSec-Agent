package checks

import (
	"sort"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// NetworkProfileTypes is a function that checks the network profile types on the system.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): An instance of RegistryKey used to access the registry keys related to network profiles.
//
// Returns:
//   - Check: A Check instance encapsulating the results of the network profile type check. The Result field of the Check instance will contain one or more of the following messages:
//   - "Network [ProfileName] is Public" if the network profile is public.
//   - "Network [ProfileName] is Private" if the network profile is private.
//   - "Network [ProfileName] is Domain" if the network profile is domain.
//   - "No network profiles found" if no network profiles are found.
//
// This function is primarily used to identify potential security risks associated with different types of network profiles on the system.
func NetworkProfileTypes(registryKey mocking.RegistryKey) Check {
	var err error
	var profilesKey mocking.RegistryKey
	var networkHashes []string
	profilesKey, err = mocking.OpenRegistryKey(registryKey, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\NetworkList\Profiles`)
	if err != nil {
		return NewCheckError(NetworkProfileTypeID, err)
	}
	defer mocking.CloseRegistryKey(profilesKey)
	networkHashes, err = profilesKey.ReadSubKeyNames(-1)
	if err != nil {
		return NewCheckErrorf(NetworkProfileTypeID, "error reading sub key names", err)
	}
	if len(networkHashes) == 0 {
		return NewCheckResult(NetworkProfileTypeID, 0, "No network profiles found")
	}
	networkNames := make(map[string]string)
	var name string
	var value []byte
	var networkKey mocking.RegistryKey
	for _, networkHash := range networkHashes {
		networkKey, err = mocking.OpenRegistryKey(profilesKey, networkHash)
		if err != nil {
			return NewCheckErrorf(NetworkProfileTypeID, "error opening network key", err)
		}
		defer mocking.CloseRegistryKey(networkKey)
		name, _, err = networkKey.GetStringValue("ProfileName")
		if err != nil {
			return NewCheckErrorf(NetworkProfileTypeID, "error reading profile name", err)
		}
		value, _, err = networkKey.GetBinaryValue("Category")
		if err != nil {
			return NewCheckErrorf(NetworkProfileTypeID, "error reading category", err)
		}
		if networkNames[name] == "" {
			switch value[0] {
			case 0:
				networkNames[name] = "Public"
			case 1:
				networkNames[name] = "Private"
			case 2:
				networkNames[name] = "Domain"
			default:
				networkNames[name] = "Unknown"
				// logger.Log.Info("Unknown network profile type")
				// logger disabled as it was crashing the tests. FIXME: ask kevin how to fix
			}
		} else {
			logger.Log.Info("Network profile " + name + " already exists")
		}
	}
	var results []string
	// Create a slice to store the network names
	keys := make([]string, 0, len(networkNames))
	// Extract the keys into the slice
	for k := range networkNames {
		keys = append(keys, k)
	}
	// Sort the slice
	sort.Strings(keys)
	// Iterate over the sorted slice and get the corresponding value from the map
	for _, k := range keys {
		results = append(results, "Network "+k+" is "+networkNames[k])
	}
	return NewCheckResult(NetworkProfileTypeID, 1, results...)
}