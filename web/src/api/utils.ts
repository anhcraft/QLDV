export default {
    parseContestSession(s: any) {
        s.questionSheet = JSON.parse(s.questionSheet)
        s.answerSheet = JSON.parse(s.answerSheet)
        if(s.hasOwnProperty("expectedAnswerSheet")) {
            s.expectedAnswerSheet = JSON.parse(s.expectedAnswerSheet)
        }
        if(s.hasOwnProperty("answerAccuracy")) {
            s.answerAccuracy = JSON.parse(s.answerAccuracy)
        }
        return s
    }
}
