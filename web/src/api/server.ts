import conf from "../conf";
import ky from 'ky';

const server = {
    loadProfile: async function (token: string) {
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
        return ky.post(`${conf.server}/profile`, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'token': token
            }
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
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
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
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
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
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
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
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
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
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
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
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
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
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
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
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
        return ky.post(`${conf.server}/get-user-stats`, {
            method: 'post',
            headers: {
                'content-type': undefined,
                'token': token
            }
        }).json();
    },
    saveProgressionChanges(data: {}, user: string, token: string) {
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
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
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
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
        if(token == null || token.length == 0) {
            return {
                "error": "CLIENT"
            };
        }
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
}

export default server;
