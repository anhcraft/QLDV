const list = {
    "USER_ILLEGAL_EMAIL": "Vui lòng dùng email của nhà trường!",
    "LOGIN_FAILED": "Đăng nhập thất bại"
}

function lookupErrorCode(s: string): string {
    if(s === undefined) return "-"
    // @ts-ignore
    return list.hasOwnProperty(s) ? list[s] : ("Lỗi: " + s)
}

export class ClientError {
    constructor(code: string) {
        this._code = code;
        this._message = lookupErrorCode(code);
    }

    get code(): string {
        return this._code;
    }

    get message(): string {
        return this._message;
    }

    private readonly _code: string;
    private readonly _message: string;
}
