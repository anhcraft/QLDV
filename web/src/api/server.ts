import conf from "../conf";
import ky from 'ky';

const server = {
    loadProfile: async function (user: string, token: string) {
        return ky.post(`${conf.server}/profile`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                'user': user
            })
        }).json();
    },
    loadPosts: function (limit: number, filterHashtags: string[], sortBy: string, lowerThan: number, belowId: any, token: string) {
        return ky.post(`${conf.server}/posts`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                'limit': limit,
                'filter_hashtags': filterHashtags,
                'below_id': parseInt(belowId),
                'sort_by': sortBy,
                'lower_than': lowerThan,
            })
        }).json();
    },
    loadPost: function (id: number, token: string) {
        return ky.get(`${conf.server}/post?id=${id}`, {
            method: 'get',
            headers: {
                'content-type': 'application/json',
                'token': token
            }
        }).json();
    },
    loadProgression: async function (token: string, user: string) {
        return ky.post(`${conf.server}/progression`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                'user': user
            })
        }).json();
    },
    changePost: async function(id: any, title: string, content: string, privacy: number, hashtag: string, removeAttachments: string[], token: string) {
        return ky.post(`${conf.server}/change-post`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                'id': id === undefined ? -1 : parseInt(id),
                'title': title,
                'hashtag': hashtag,
                'content': content,
                'privacy': privacy,
                'remove_attachments': removeAttachments
            })
        }).json();
    },
    removePost: async function(id: number, token: string) {
        return ky.post(`${conf.server}/remove-post`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token,
                'id': id.toString()
            }
        }).json();
    },
    updatePostStat: async function(id: number, action: string, token: string) {
        return ky.post(`${conf.server}/update-post-stat`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token,
                'id': id.toString(),
                'action': action
            }
        }).json();
    },
    uploadPostAttachment: async function(id: number, attachment: Blob, token: string) {
        return ky.post(`${conf.server}/upload-attachment`, {
            method: 'post',
            headers: {
                'content-type': undefined,
                'token': token,
                'id': id.toString()
            },
            body: attachment
        }).json();
    },
    loadUsers: async function (limit: number, offset: number, filter: { name: string; class: string; email: string; certified: number }, token: string) {
        return ky.post(`${conf.server}/users`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                'limit': limit,
                'offset': offset,
                'filter_name': filter.name,
                'filter_class': filter.class,
                'filter_email': filter.email,
                'filter_certified': filter.certified
            })
        }).json();
    },
    saveUserChanges: async function(changes: {}, token: string) {
        return ky.post(`${conf.server}/change-users`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify(changes)
        }).json();
    },
    getUserStats: async function(token: string) {
        return ky.post(`${conf.server}/get-user-stats`, {
            method: 'post',
            headers: {
                'content-type': undefined,
                'token': token
            }
        }).json();
    },
    saveProgressionChanges(data: {}, user: string, token: string) {
        return ky.post(`${conf.server}/change-progression`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify(Object.assign(data, {
                user: user
            }))
        }).json();
    },
    loadEvents: function (limit: number, belowId: number, beginDate: number, endDate: number, token: string) {
        return ky.get(`${conf.server}/events?limit=${limit}&below-id=${belowId}&begin-date=${beginDate}&end-date=${endDate}`, {
            method: 'get',
            headers: {
                'content-type': 'application/json',
                'token': token
            }
        }).json();
    },
    removeEvent: async function(id: number, token: string) {
        return ky.post(`${conf.server}/remove-event`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token,
                'id': id.toString()
            }
        }).json();
    },
    changeEvent: async function (id: any, event: { title: string, beginDate: Date, endDate: Date, privacy: number }, token: string) {
        return ky.post(`${conf.server}/change-event`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                'id': parseInt(id),
                'title': event.title,
                'begin_date': event.beginDate.getTime(),
                'end_date': event.endDate.getTime(),
                'privacy': event.privacy
            })
        }).json();
    },
    loadEvent: function (id: number, token: string) {
        return ky.get(`${conf.server}/event?id=${id}`, {
            method: 'get',
            headers: {
                'content-type': 'application/json',
                'token': token
            }
        }).json();
    },
    setProfileCover: async function(file: Blob, token: string) {
        return ky.post(`${conf.server}/set-profile-cover`, {
            method: 'post',
            headers: {
                'content-type': undefined,
                'token': token
            },
            body: file
        }).json();
    },
    setProfileBoard: async function(board: string, token: string) {
        return ky.post(`${conf.server}/set-profile-board`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                "board": board
            })
        }).json();
    },
    setProfileSettings: async function(settings: number, token: string) {
        return ky.post(`${conf.server}/set-profile-settings`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                "settings": settings
            })
        }).json();
    },
    changeContest(id: any, contest: { limitTime: number; limitQuestions: number; limitSessions: number; dataSheet: []; acceptingAnswers: boolean, info: string }, token: string) {
        return ky.post(`${conf.server}/change-contest`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                id: parseInt(id),
                accepting_answers: contest.acceptingAnswers,
                limit_questions: contest.limitQuestions,
                limit_sessions: contest.limitSessions,
                limit_time: contest.limitTime * 60000,
                data_sheet: JSON.stringify(contest.dataSheet),
                info: contest.info
            })
        }).json();
    },
    removeContest(id: number, token: string) {
        return ky.post(`${conf.server}/remove-contest`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token,
                'id': id.toString()
            }
        }).json();
    },
    submitContestSession(id: string, answerSheet: string, saveOnly: boolean, token: string) {
        return ky.post(`${conf.server}/submit-contest-session`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                'id': id,
                'answer_sheet': JSON.stringify(answerSheet),
                'save_only': saveOnly
            })
        }).json();
    },
    joinContestSession(id: string, token: string) {
        return ky.post(`${conf.server}/join-contest-session`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                id: id
            })
        }).json();
    },
    loadContestSessions: function (contest: any, limit: number, offset: number, filterAttendant: string, filterFinished: boolean, sortBy: string[], token: string) {
        return ky.post(`${conf.server}/contest-sessions`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                'contest': parseInt(contest),
                'limit': limit,
                'offset': offset,
                'filter_attendant': filterAttendant,
                'filter_finished': filterFinished,
                'sort_by': sortBy
            })
        }).json();
    },
    getContestStats: async function(contest: any, token: string) {
        return ky.post(`${conf.server}/get-contest-stats`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                'id': parseInt(contest)
            })
        }).json();
    },
    getHashtags: async function() {
        return ky.get(`${conf.server}/get-hashtags`, {
            method: 'get',
            headers: {
                'content-type': 'application/json'
            }
        }).json().then((t: any) => t.hashtags);
    }
}

export default server;
