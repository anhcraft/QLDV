import Cookies from 'js-cookie';

const auth = {
    getToken: function (): string | undefined {
        return Cookies.get('qldvtkn')
    },
    destroySession: function () {
        Cookies.remove('qldvtkn');
    }
};

export default auth;
