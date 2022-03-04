const list = {
    "ERR_UNKNOWN_USER": "Tài khoản không tồn tại",
    "ERR_TOKEN_VERIFY": "Vui lòng đăng nhập!",
    "ERR_NO_PERMISSION": "Vui lòng đăng nhập!",
    "ERR_UNKNOWN_EVENT": "Không tìm thấy sự kiện",
    "ERR_INVALID_EVENT_ID": "Không tìm thấy sự kiện",
    "ERR_EVENT_TITLE_MIN": "Tiêu đề sự kiện phải từ 5 - 300 kí tự",
    "ERR_EVENT_TITLE_MAX": "Tiêu đề sự kiện phải từ 5 - 300 kí tự",
    "ERR_EVENT_UNAVAILABLE": "Không thể tham gia sự kiện ngay lúc này",
    "ERR_UNKNOWN_CONTEST": "Không tìm thấy cuộc thi",
    "ERR_INVALID_CONTEST_ID": "Không tìm thấy cuộc thi",
    "ERR_CONTEST_DATA_INSUFFICIENT": "Lỗi cuộc thi. Vui lòng báo quản trị viên!",
    "ERR_CONTEST_CLOSED": "Cuộc thi đã kết thúc",
    "ERR_CONTEST_ATTENDED": "Đã tham gia cuộc thi",
    "ERR_CONTEST_ATTENDED_MAX": "Số lần tham gia đã đạt giới hạn",
    "ERR_DATE_RANGE": "Ngày kết thúc phải sau ngày bắt đầu",
    "ERR_UNKNOWN_POST": "Không tìm thấy bài viết",
    "ERR_POST_TITLE_MIN": "Tiêu đề bài viết phải từ 5 - 300 kí tự",
    "ERR_POST_TITLE_MAX": "Tiêu đề bài viết phải từ 5 - 300 kí tự",
    "ERR_POST_CONTENT_MIN": "Nội dung bài viết phải từ 10 - 100.000 kí tự",
    "ERR_POST_CONTENT_MAX": "Nội dung bài viết phải từ 10 - 100.000 kí tự",
    "ERR_ILLEGAL_ATTACHMENT": "Ảnh đính kèm phải ở định dạng .png/.jpg/.jpeg",
    "ERR_ATTACHMENT_REMOVE_FAILED": "Lỗi xóa ảnh đính kèm",
    "ERR_UNKNOWN_POST_ACTION": "Hành động tương tác bài việt không hợp lệ",
    "ERR_ILLEGAL_PROFILE_COVER": "Ảnh bìa phải ở định dạng .png/.jpg/.jpeg",
    "ERR_PROFILE_BOARD_CONTENT_MIN": "Nội dung tường nhà phải từ 10 - 10.000 kí tự",
    "ERR_PROFILE_BOARD_CONTENT_MAX": "Nội dung tường nhà phải từ 10 - 10.000 kí tự"
}

export default function lookupErrorCode(s: string) {
    if(s === undefined) return "-"
    // @ts-ignore
    return list.hasOwnProperty(s) ? list[s] : ("Mã lỗi: " + s)
}
