const RoleGuest: number = 0
const RoleRegularMember: number = 1
const RoleCertifiedMember: number = 2
const RoleClassDeputySecretary: number = 3
const RoleClassSecretary: number = 4
const RoleDeputySecretary: number = 5
const RoleSecretary: number = 6

const RoleGroupGuest: number = 0
const RoleGroupMember: number = 1
const RoleGroupClassManager: number = 2
const RoleGroupGlobalManager: number = 3

export function GetRoleName(role: number): string {
    switch (role) {
        case RoleRegularMember:
            return "Đoàn viên (Chưa gia nhập)"
        case RoleCertifiedMember:
            return "Đoàn viên"
        case RoleClassDeputySecretary:
            return "Phó bí thư Chi Đoàn"
        case RoleClassSecretary:
            return "Bí thư Chi Đoàn"
        case RoleDeputySecretary:
            return "Phó bí thư Đoàn"
        case RoleSecretary:
            return "Bí thư Đoàn"
    }
    return "Khách"
}

export function GetRoleGroup(role: number): number {
    switch (role) {
        case RoleGuest:
            return RoleGroupGuest
        case RoleCertifiedMember:
        case RoleRegularMember:
            return RoleGroupMember
        case RoleClassSecretary:
        case RoleClassDeputySecretary:
            return RoleGroupClassManager
        case RoleSecretary:
        case RoleDeputySecretary:
            return RoleGroupGlobalManager
    }
    return RoleGuest
}

export function IsMember(role: number): boolean {
    return role == RoleRegularMember || role == RoleCertifiedMember
}

export function IsManager(role: number): boolean {
    return GetRoleGroup(role) == RoleGroupClassManager || GetRoleGroup(role) == RoleGroupGlobalManager
}

export function IsClassRole(role: number): boolean {
    return IsMember(role) || GetRoleGroup(role) == RoleGroupClassManager
}

export function IsLoggedIn(role: number): boolean {
    return role != RoleGuest
}

export function IsCertified(role: number): boolean {
    return role != RoleGuest && role != RoleRegularMember
}
