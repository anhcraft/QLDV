import conf from "../conf";
import auth from "./auth";
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
    loadPosts: function (limit: number, olderThan: number) {
        return ky.get(`${conf.server}/posts?limit=${limit}&older=${olderThan}`, {
            method: 'get',
            headers: {
                'content-type': 'application/json'
            }
        }).json();
    },
    loadPost: function (id: string) {
        return ky.get(`${conf.server}/post?id=${id}`, {
            method: 'get',
            headers: {
                'content-type': 'application/json'
            }
        }).json();
    },
    loadProgression: async function (token: string) {
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
            }
        }).json();
    },
    changePost: async function(id: string, title: string, content: string, token: string) {
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
                'content': content
            })
        }).json();
    }
}

export default server;