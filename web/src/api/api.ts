import ky from "ky";
import conf from "../conf";
import auth from "../auth/auth";
import {ServerError} from "./server-error";

const list = {
    //"ERROR_INVALID_REQUEST_QUERY": "",
    //"ERROR_INVALID_REQUEST_BODY": "",
    "ERROR_TOKEN_VERIFY": "Không thể xác thực người dùng",
    "ERROR_INVALID_TOKEN": "Sai mã xác thực",
    "ERROR_UNKNOWN_TOKEN_OWNER": "Không thể tìm thấy người dùng",
    "ERROR_NO_PERMISSION": "Không có quyền truy cập",
    "ERROR_UNKNOWN_USER": "Không tìm thấy người dùng",
    "ERROR_PROFILE_BOARD_TOO_SHORT": "Bảng tin quá ngắn (10-10.000 kí tự)",
    "ERROR_PROFILE_BOARD_TOO_LONG": "Bảng tin quá dài (10-10.000 kí tự)",
    "ERROR_PROFILE_COVER_TOO_LARGE": "Kích thước ảnh đại diện quá lớn (Tối đa 500KB)",
    "ERROR_PROFILE_COVER_UPLOAD_FAILED": "Tải ảnh đại diện thất bại",
    "ERROR_UNSUPPORTED_PROFILE_COVER": "Ảnh đại diện không phù hợp",
    "ERROR_UNKNOWN_POST": "Không tìm thấy bài viết",
    "ERROR_POST_TITLE_TOO_SHORT": "Tiêu đề bài viết quá ngắn (10-300 kí tự)",
    "ERROR_POST_TITLE_TOO_LONG": "Tiêu đề bài viết quá dài (10-300 kí tự)",
    "ERROR_POST_HEADLINE_TOO_LONG": "Đề mục bài viết quá dài (30-250 kí tự)",
    "ERROR_POST_HEADLINE_TOO_SHORT": "Đề mục bài viết quá ngắn (30-250 kí tự)",
    "ERROR_POST_CONTENT_TOO_LONG": "Nội dung bài viết quá dài (100-100.000 kí tự)",
    "ERROR_POST_CONTENT_TOO_SHORT": "Nội dung bài viết quá ngắn (100-100.000 kí tự)",
    "ERROR_POST_HASHTAG_TOO_LONG": "Hashtag bài viết quá dài (5-20 kí tự)",
    "ERROR_POST_HASHTAG_TOO_SHORT": "Hashtag bài viết quá ngắn (5-20 kí tự)",
    "ERROR_INVALID_POST_HASHTAG": "Hashtag bài viết không đúng định dạng",
    "ERROR_POST_CREATE_FAILED": "Tạo bài viết thất bại",
    "ERROR_POST_UPDATE_FAILED": "Cập nhật bài viết thất bại",
    "ERROR_POST_DELETE_FAILED": "Xóa bài viết thất bại",
    "ERROR_POST_STAT_UPDATE_FAILED": "Cập nhật thống kê bài viết thất bại",
    "ERROR_POST_HASHTAG_LIST_FAILED": "Lấy danh sách hashtag thất bại",
    "ERROR_ATTACHMENT_TOO_LARGE": "Ảnh đính kèm quá lớn (Tối đa 2MB)",
    "ERROR_UNSUPPORTED_ATTACHMENT": "Định dạng ảnh đính kèm không được hỗ trợ",
    "ERROR_ATTACHMENT_DELETE_FAILED": "Xóa ảnh đính kèm thất bại",
    "ERROR_UNKNOWN_EVENT": "Không tìm thấy sự kiện",
    "ERROR_EVENT_TITLE_TOO_SHORT": "Tiêu đề sự kiện quá ngắn (10-300 kí tự)",
    "ERROR_EVENT_TITLE_TOO_LONG": "Tiêu đề sự kiện quá dài (10-300 kí tự)",
    "ERROR_EVENT_CREATE_FAILED": "Tạo sự kiện thất bại",
    "ERROR_EVENT_UPDATE_FAILED": "Cập nhật sự kiện thất bại",
    "ERROR_EVENT_DELETE_FAILED": "Xóa sự kiện thất bại",
}

const API = {
    lookupErrorCode(s: string): string {
        if(s === undefined) return "-"
        // @ts-ignore
        return list.hasOwnProperty(s) ? list[s] : ("Lỗi: " + s)
    },
    postObject(endpoint: string, data: any): Promise<any | ServerError> {
        return ky.post(conf.server + endpoint, {
            method: 'post',
            headers: {
                'content-type': 'application/json',
                'access-token': auth.getToken()
            },
            body: JSON.stringify(data)
        }).json().then((res: any) => {
            if(res["success"]) {
                return res["result"]
            } else {
                return new ServerError(res["error"])
            }
        });
    },
    getObject(endpoint: string, data: any): Promise<any | ServerError> {
        return ky.post(conf.server + endpoint, {
            method: 'get',
            headers: {
                'content-type': 'application/json',
                'access-token': auth.getToken()
            },
            body: JSON.stringify(data)
        }).json().then((res: any) => {
            if(res["success"]) {
                return res["result"]
            } else {
                return new ServerError(res["error"])
            }
        });
    }
}

export default API;
