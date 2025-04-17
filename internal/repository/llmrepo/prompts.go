package llmrepo

import (
	"fmt"
)

// store system prompts here
// This is strait form Chat-GTP, still need assessed.
func (o *OllamaRepo) JobListingPrompt(jobDescription string) string {
	return fmt.Sprintf(`
        You are an AI assistant extracting job details. 
        Analyze the following job description and return a JSON object with the following fileds. 
		If data is not present for any field, return an empty string for the value in that field:
		- "job_id": extracted joblisting ID, Requisition ID Job# or similar.
		- Return the company as a JSON object with "company_name".
		Example:
		{
		"company": {
			"company_name": "Microsoft"
			}
		}
		Theses Examples are incorrect:

		{
		"company_name": "Microsoft"
		}
		{
		"company": "Microsoft"
		}
		- "job_title": The job title
		- "job_description": summarize the job description into one or two short paragraphs
		- "req_yoe":  The years of experience required. This should be a single number. If there is a range, use the 
		smallest number. Only use a single interger as the value!
		- "req_skills": A list of required soft and hard skills. Extract the required skills from the job listing 
		and return them as a JSON array. Each skill should be a short, precise phrase describing a technical or 
		professional ability. Do not include things like "Hands-on experience with" or "demonstrated experience in".
        - "req_certs": A list of required certificatios. Extract the required skills from the job listing 
		and return them as a JSON array. Each certification should be a short, precise acronym or disignation. 
		If a security clearance requirement is listed, add it in to this section as a certification, i.e DOEQ, TS/SCI
		- "low_pay": If a pay range is is listed, this will be the lower of the two. If a single value is listed it will go here. 
		- "high_pay": If a pay range is is listed, this will be the higher of the two. If a single value is listed leave this value as an empty string. 
		Return the location as a JSON object with "city" and "state". This values is the state in a two letter code if listed. 
		Example:
		{
		"location": {
			"city": "San Diego",
			"state": "CA",
		}
		}
		This Example is incorrect:
		{
		"location": "San Diego, CA"
		}
        Do not include a pay_range field, only put the low and high pay if they are listed. Respond only in JSON format. only include those fields listed above in one set of currly braces. 
		Here is an example format for the full response:
		{
		"job_id": "R-00153190",
		"company": {
			"company_name": "Microsoft"
			},
		"job_title": "Cybersecurity Engineer",
		"job_description": "summary of job description",
		"req_yoe": 5,
		"req_skills": ["Network Defense", "PCAP Analysis", "Computer Incident Response"],
		"req_certs": ["CompTIA CYSA+", "CISSP", "TS/SCI"],
		"low_pay": 90000,
		"high_pay": 120000,
		"location": {
			"city": "Denver",
			"state": "CO",
			}
		}
        Job Description:
        %s
    `, jobDescription)
}
