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
    loadPosts: function (limit: number, olderThan: number, token: string) {
        return ky.get(`${conf.server}/posts?limit=${limit}&older=${olderThan}`, {
            method: 'get',
            headers: {
                'content-type': 'application/json',
                'token': token
            }
        }).json();
    },
    loadPost: function (id: string, token: string) {
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
    changePost: async function(id: string, title: string, content: string, privacy: number, removeAttachments: string[], token: string) {
        return ky.post(`${conf.server}/change-post`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                'id': id === undefined ? '' : id,
                'title': title,
                'content': content,
                'privacy': privacy,
                'remove_attachments': removeAttachments
            })
        }).json();
    },
    removePost: async function(id: string, token: string) {
        return ky.post(`${conf.server}/remove-post`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token,
                'id': id === undefined ? '' : id
            }
        }).json();
    },
    updatePostStat: async function(id: string, action: string, token: string) {
        return ky.post(`${conf.server}/update-post-stat`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token,
                'id': id === undefined ? '' : id,
                'action': action
            }
        }).json();
    },
    uploadPostAttachment: async function(id: string, attachment: Blob, token: string) {
        return ky.post(`${conf.server}/upload-attachment`, {
            method: 'post',
            headers: {
                'content-type': undefined,
                'token': token,
                'id': id
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
    loadEvents: function (limit: number, olderThan: number, fromDate: number, toDate: number, token: string) {
        return ky.get(`${conf.server}/events?limit=${limit}&older=${olderThan}&from-date=${fromDate}&to-date=${toDate}`, {
            method: 'get',
            headers: {
                'content-type': 'application/json',
                'token': token
            }
        }).json();
    },
    removeEvent: async function(id: string, token: string) {
        return ky.post(`${conf.server}/remove-event`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token,
                'id': id === undefined ? '' : id
            }
        }).json();
    },
    changeEvent: async function (id: string, event: { endDate: Date, title: string, startDate: Date, privacy: number }, token: string) {
        return ky.post(`${conf.server}/change-event`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                'id': id === undefined ? '' : id,
                'title': event.title,
                'start_date': event.startDate.getTime(),
                'end_date': event.endDate.getTime(),
                'privacy': event.privacy
            })
        }).json();
    },
    loadEvent: function (id: string, token: string) {
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
    changeContest(id: string, contest: { limitTime: number; limitQuestions: number; dataSheet: []; acceptingAnswers: boolean, info: string }, token: string) {
        return ky.post(`${conf.server}/change-contest`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            },
            body: JSON.stringify({
                id: id,
                accepting_answers: contest.acceptingAnswers,
                limit_questions: contest.limitQuestions,
                limit_time: contest.limitTime * 60000,
                data_sheet: JSON.stringify(contest.dataSheet),
                info: contest.info
            })
        }).json();
    },
    removeContest(id: string, token: string) {
        return ky.post(`${conf.server}/remove-contest`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token,
                'id': id === undefined ? '' : id
            }
        }).json();
    },
    loadContestSession: function (id: string, token: string) {
        return ky.get(`${conf.server}/contest-session?id=${id}`, {
            method: 'get',
            headers: {
                'content-type': 'application/json',
                'token': token
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
                id: id,
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
    loadContestSessions: function (contest: string, limit: number, olderThan: number, token: string) {
        return ky.get(`${conf.server}/contest-sessions?contest=${contest}&limit=${limit}&older=${olderThan}`, {
            method: 'get',
            headers: {
                'content-type': 'application/json',
                'token': token
            }
        }).json();
    },
}

export default server;
