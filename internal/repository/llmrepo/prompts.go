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
		- "job_title": The job title
		- "job_description": summarize the job description into one or two short paragraphs
		- "req_yoe": The years of experience required. This should be a single number. If there is a range, use the smallest number.
		- "req_skills": A list of required soft and hard skills. Use one or two words for each skill and seperate several skills with a comma. Keep these as short as possible, do not include things like "Hands-on experience with" or "demonstrated experience in".
        - "req_certs": A list of required certificatios. seperate each with a comma. 
		- "low_pay": If a pay range is is listed, this will be the lower of the two. If a single value is listed it will go here. 
		- "high_pay": If a pay range is is listed, this will be the higher of the two. If a single value is listed leave this value as an empty string. 
		- "location_city": This values is the city if listed.
		- "location_state": This values is the state in a two letter code if listed. 

        Do not include a pay_range field, only put the low and high pay if they are listed. Respond only in JSON format. only include those fields listed above in one set of currly braces. 

        Job Description:
        %s
    `, jobDescription)
}
