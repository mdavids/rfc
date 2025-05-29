package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
)

// showForSaleRecords retrieves and displays ONLY specific '_for-sale' TXT records with an fcod
// starting with "NLFS-" followed by 48 characters, based on
// draft-davids-forsalereg-08.txt.
// All other record types or malformed fcod values are silently skipped.
func showForSaleRecords(domain string) {
	forSaleSubdomain := fmt.Sprintf("_for-sale.%s", domain)
	fmt.Printf("Attempting to retrieve and display specific '_for-sale' records for: %s\n\n", forSaleSubdomain)

	txtRecords, err := net.LookupTXT(forSaleSubdomain)
	if err != nil {
		if dnsErr, ok := err.(*net.DNSError); ok {
			if dnsErr.IsNotFound {
				fmt.Printf("No '_for-sale' TXT record found for %s. Domain is likely not declared for sale via this method.\n", domain)
			} else {
				fmt.Printf("An error occurred while querying DNS records: %v\n", err)
			}
		} else {
			fmt.Printf("An error occurred while querying DNS records: %v\n", err)
		}
		return
	}

	// Sort the raw records for consistent display order
	sort.Strings(txtRecords)

	foundSpecificNlfsFcod := false
	fmt.Println("--- Specific For Sale Information (NLFS- fcod) ---")
	recordCount := 0

	// Specific fcod pattern: "NLFS-" followed by exactly 48 characters
	nlfsFcodPattern := regexp.MustCompile(`^NLFS-(.{48})$`)

	// The specific 48-character string that triggers the special action
	specificIdentifierString := "NGYyYjEyZWYtZTUzYi00M2U0LTliNmYtNTcxZjBhMzA2NWQy"

	for _, txtRecordContent := range txtRecords {
		versionTag := "v=FORSALE1;"

		// Apply robustness principle: allow optional whitespace after the version tag.
		if !strings.HasPrefix(txtRecordContent, versionTag) {
			continue // Skip records not starting with the correct version tag
		}

		contentPart := strings.TrimSpace(txtRecordContent[len(versionTag):])

		// ONLY consider records that are exactly an 'fcod=' tag
		if !strings.HasPrefix(contentPart, "fcod=") {
			continue // Skip if not an fcod tag
		}

		// Extract the fcod value
		fcodValue := contentPart[len("fcod="):]

		// Now, specifically check if this fcod value matches the NLFS- pattern
		nlfsMatch := nlfsFcodPattern.FindStringSubmatch(fcodValue)

		if len(nlfsMatch) > 1 { // If nlfsMatch has submatches, it means the pattern matched
			recordCount++
			fmt.Printf("\nRecord %d:\n", recordCount)
			foundSpecificNlfsFcod = true

			identifierPart := nlfsMatch[1] // The 48 characters after "NLFS-"

			fmt.Printf("  For Sale Code (NLFS- format found): %s\n", fcodValue)

			if identifierPart == specificIdentifierString {
				fmt.Printf("  Action: Visit https://www.sidn.nl/en/landing-page-buying-and-selling-example?domain=%s for more information.\n", domain)
			} else {
				fmt.Printf("  Warning: Specific fcod not found in database.\n")
			}
		}
		// If it's an fcod but doesn't match the NLFS- pattern, it's silently skipped as per requirements.
	}

	if !foundSpecificNlfsFcod {
		fmt.Println("\nNo '_for-sale' TXT records with a specific 'NLFS-' fcod pattern were found.")
		fmt.Println("Other types of '_for-sale' records (ftxt, furi, or fcod not matching 'NLFS-{.48}')")
		fmt.Println("are ignored by this script. Use 'check_dns_for_sale.py' for full validation.")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run show_dns_for_sale.go <domain_name>")
		fmt.Println("Example: go run show_dns_for_sale.go example.com")
		os.Exit(1)
	}

	domainToShow := os.Args[1]
	showForSaleRecords(domainToShow)
}
