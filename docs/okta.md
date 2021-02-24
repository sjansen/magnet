1. Add Application
1. Create New App
 - Do not use a template
 - Platform: Web
 - Sign on method: SAML 2.0
1. SAML Settings
 - Single sign on URL: https://example.com/saml/acs
 - Use this for Recipient URL and Destination URL: true
 - Allow this app to request other SSO URLs: false
 - Audience URI (SP Entity ID): magnet
 - Default RelayState: blank
 - Name ID format: Unspecified
 - Application username: email
 - Update application username on: Create and update
1. Attribute Statements
 - firstName Unspecified user.firstName
 - lastName Unspecified user.lastName
 - email Unspecified user.email
 - username Unspecified String.substringBefore(user.email, "@") 
1. Group Attribute Statements
 - roles Unspecified Starts with: magnet- 
1. Assignments
 - Assign > Assign to People
