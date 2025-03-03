package esco

import (
	"github.com/threehook/esco-search/csv"
)

func CreateSkillsByOccupationMap(occupationSkillRelationsFilePath, skillsFilePath string) (map[string]string, error) {
	columns := []string{"occupationUri", "skillUri"}
	skillsByOccupation, err := csv.ReadCSV(occupationSkillRelationsFilePath, columns...)
	if err != nil {
		return nil, err
	}

	columns = []string{"conceptUri", "description"}

	skillDescriptionsByConceptUri, err := csv.ReadCSV(skillsFilePath, columns...)
	if err != nil {
		return nil, err
	}

	skillsMap := make(map[string]string, len(skillDescriptionsByConceptUri))
	for _, skill := range skillDescriptionsByConceptUri {
		skillsMap[skill[0]] = skill[1]
	}

	occupationsWithSkills := addSkillsToOccupationRefs(skillsByOccupation, skillsMap)

	return occupationsWithSkills, nil
}

func CreateEducationTypesLevelsByOccupationMap(educationsFilePath string) (map[string]string, error) {
	columns := []string{"code", "educationTypesAndLevels"}
	educationTypesLevelsByOccupationCode, err := csv.ReadCSV(educationsFilePath, columns...)
	if err != nil {
		return nil, err
	}

	educationTypesLevelsMap := make(map[string]string, len(educationTypesLevelsByOccupationCode))
	for _, education := range educationTypesLevelsByOccupationCode {
		educationTypesLevelsMap[education[0]] = education[1]
	}

	return educationTypesLevelsMap, nil
}

func addSkillsToOccupationRefs(skillsByOccupation [][]string, skillsMap map[string]string) map[string]string {
	skillsByOccupationRefMap := make(map[string]string, 0)
	for _, occupation := range skillsByOccupation {
		// Find the skill description
		description := skillsMap[occupation[1]]

		// Find the occupation and add the new skill
		if value, ok := skillsByOccupationRefMap[occupation[0]]; ok {
			desc := skillsByOccupationRefMap[value]
			skillsByOccupationRefMap[occupation[0]] = desc + description
		} else {
			skillsByOccupationRefMap[occupation[0]] = description
		}
	}

	return skillsByOccupationRefMap
}

func AddSkillsToOccupations(occupations [][]string, skillsByOccupationRefMap map[string]string) {
	for _, occupation := range occupations {
		skills := skillsByOccupationRefMap[occupation[0]]
		occupation[0] = skills
	}
}

func AddEducationTypesAndLevelsToOccupations(occupations [][]string, educationByOccupationRefMap map[string]string) {
	for i, occupation := range occupations {
		occupations[i] = append(occupations[i], educationByOccupationRefMap[occupation[5]])
	}
}
