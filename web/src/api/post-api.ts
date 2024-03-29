import API from "./api";
import {ServerError} from "./server-error";

const PostAPI = {
    getPost(id: string): Promise<any | ServerError> {
        return API.getObject("/post/" + id, undefined)
    },
    updatePost(id: string, data: any): Promise<any | ServerError> {
        return API.postObject("/post/" + id, data)
    },
    deletePost(id: string): Promise<any | ServerError> {
        return API.deleteObject("/post/" + id, {})
    },
    listPosts(queries: {
        "limit": number,
        "filter-hashtags": string[],
        "sort-by": "like" | "view" | "date",
        "lower-than": number,
        "below-id": number
    }): Promise<any[] | ServerError> {
        return API.getObject("/posts/", queries).then(v => {
            return v instanceof ServerError ? v : v["posts"]
        })
    },
    updatePostStat(id: string, data: {
        "like": boolean | undefined,
        "view": boolean | undefined,
    }): Promise<any | ServerError> {
        return API.postObject("/post-stat/" + id, data)
    },
    uploadAttachment(postId: string, data: Blob): Promise<any | ServerError> {
        return API.postBlob("/post-attachment/" + postId, data)
    },
    deleteAttachment(attId: string[]): Promise<any | ServerError> {
        return API.deleteObject("/post-attachment/", {
            "id": attId
        }).then(v => {
            return v instanceof ServerError ? v : v["remaining"]
        })
    },
    getHashtags(): Promise<any | ServerError> {
        return API.getObject("/post-hashtags/", {}).then(v => {
            return v instanceof ServerError ? v : v["hashtags"]
        })
    }
}

export default PostAPI;
