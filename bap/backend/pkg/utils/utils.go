package utils

import (
	"fmt"
	"strings"
)

// GetStateCode returns the code for a given Indian state.
func GetStateCode(state string) (string, error) {
	stateCodes := map[string]string{
		"Andhra Pradesh":              "IN-AP",
		"Arunachal Pradesh":           "IN-AR",
		"Assam":                       "IN-AS",
		"Bihar":                       "IN-BR",
		"Chhattisgarh":                "IN-CG",
		"Goa":                         "IN-GA",
		"Gujarat":                     "IN-GJ",
		"Haryana":                     "IN-HR",
		"Himachal Pradesh":            "IN-HP",
		"Jharkhand":                   "IN-JH",
		"Karnataka":                   "IN-KA",
		"Kerala":                      "IN-KL",
		"Madhya Pradesh":              "IN-MP",
		"Maharashtra":                 "IN-MH",
		"Manipur":                     "IN-MN",
		"Meghalaya":                   "IN-ML",
		"Mizoram":                     "IN-MZ",
		"Nagaland":                    "IN-NL",
		"Odisha":                      "IN-OD",
		"Punjab":                      "IN-PB",
		"Rajasthan":                   "IN-RJ",
		"Sikkim":                      "IN-SK",
		"Tamil Nadu":                  "IN-TN",
		"Telangana":                   "IN-TS",
		"Tripura":                     "IN-TR",
		"Uttar Pradesh":               "IN-UP",
		"Uttarakhand":                 "IN-UK",
		"West Bengal":                 "IN-WB",
		"Andaman and Nicobar Islands": "IN-AN",
		"Dadra and Nagar Haveli and Daman and Diu": "IN-DN",
		"Delhi":             "IN-DL",
		"Lakshadweep":       "IN-LD",
		"Puducherry":        "IN-PY",
		"Pondicherry":       "IN-PY",
		"Ladakh":            "IN-LA",
		"Jammu and Kashmir": "IN-JK",
	}

	if code, exists := stateCodes[state]; exists {
		return code, nil
	}

	return "", fmt.Errorf("state code not found for %s", state)
}

// GetCityCode returns a three-letter code for a given Indian city.
func GetCityCode(city string) (string, error) {
	cityCodes := map[string]string{
		"New Delhi":          "NDL",
		"Mumbai":             "MUM",
		"Bangalore":          "BLR",
		"Hyderabad":          "HYD",
		"Chennai":            "CHE",
		"Kolkata":            "KOL",
		"Pune":               "PUN",
		"Ahmedabad":          "AMD",
		"Jaipur":             "JAI",
		"Lucknow":            "LKO",
		"Chandigarh":         "CHD",
		"Patna":              "PAT",
		"Bhopal":             "BPL",
		"Thiruvananthapuram": "TVM",
		"Vadodara":           "VAD",
		"Nagpur":             "NAG",
		"Indore":             "IND",
		"Visakhapatnam":      "VIZ",
		"Surat":              "SUR",
		"Coimbatore":         "CBE",
		"Kanpur":             "KNP",
		"Agra":               "AGR",
		"Varanasi":           "VNS",
		"Ranchi":             "RNC",
		"Jamshedpur":         "JSP",
		"Dehradun":           "DED",
		"Shimla":             "SHL",
		"Srinagar":           "SXR",
		"Guwahati":           "GUW",
		"Shillong":           "SHG",
		"Imphal":             "IMP",
		"Aizawl":             "AIZ",
		"Kohima":             "KOH",
		"Itanagar":           "ITA",
		"Panaji":             "PNJ",
		"Raipur":             "RPR",
		"Amritsar":           "ASR",
		"Ludhiana":           "LDH",
		"Faridabad":          "FRD",
		"Ghaziabad":          "GZB",
		"Noida":              "NOD",
		"Gurgaon":            "GGN",
		"Madurai":            "MDU",
		"Trichy":             "TRZ",
		"Salem":              "SLM",
		"Jodhpur":            "JOD",
		"Udaipur":            "UDA",
		"Kota":               "KOT",
		"Meerut":             "MRT",
		"Aligarh":            "ALG",
		"Gwalior":            "GWL",
		"Jabalpur":           "JBP",
		"Ajmer":              "AJM",
		"Nashik":             "NSK",
		"Kolhapur":           "KOP",
		"Solapur":            "SLP",
		"Hubli":              "HBL",
		"Belgaum":            "BEL",
		"Mysore":             "MYS",
		"Rajkot":             "RAJ",
		"Jamnagar":           "JAM",
		"Bhavnagar":          "BHV",
		"Anand":              "AND",
		"Bhuj":               "BHJ",
		"Porbandar":          "PBD",
		"Jalandhar":          "JAL",
		"Ambala":             "AMB",
		"Bilaspur":           "BSP",
		"Durgapur":           "DGP",
		"Asansol":            "ASN",
		"Howrah":             "HWH",
	}

	// Normalize input city name to handle case sensitivity and trim spaces
	if code, exists := cityCodes[strings.TrimSpace(city)]; exists {
		return code, nil
	}

	return "", fmt.Errorf("city code not found for %s", city)
}

func GetCountryCode(country string) (string, error) {
	countryCodes := map[string]string{
		"India":          "IND",
	}

	// Normalize input country name to handle case sensitivity and trim spaces
	if code, exists := countryCodes[strings.TrimSpace(country)]; exists {
		return code, nil
	}

	return "", fmt.Errorf("country code not found for %s", country)
}

func GetLanguageCode(language string) (string, error) {
	languageCodes := map[string]string{
		"English":          "en",
		"Hindi":            "hi",
		"Kannada":          "kn",
	}

	// Normalize input language name to handle case sensitivity and trim spaces
	if code, exists := languageCodes[strings.TrimSpace(language)]; exists {
		return code, nil
	}

	return "", fmt.Errorf("language code not found for %s", language)
}

