export const RoleGuest: number = 0
export const RoleRegularMember: number = 1
export const RoleCertifiedMember: number = 2
export const RoleClassDeputySecretary: number = 3
export const RoleClassSecretary: number = 4
export const RoleDeputySecretary: number = 5
export const RoleSecretary: number = 6
export const RoleRoot: number = 7

export const RoleGroupGuest: number = 0
export const RoleGroupMember: number = 1
export const RoleGroupClassManager: number = 2
export const RoleGroupGlobalManager: number = 3
export const RoleGroupRoot: number = 4

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
        case RoleRoot:
            return "-"
    }
    return "Khách"
}

export function GetRoleTable(): { role: number; name: string }[] {
    return [
        RoleGuest,
        RoleRegularMember,
        RoleCertifiedMember,
        RoleClassDeputySecretary,
        RoleClassSecretary,
        RoleDeputySecretary,
        RoleSecretary,
        //RoleRoot, // hidden
    ].map(v => {
        return {
            "role": v,
            "name": GetRoleName(v)
        }
    })
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
        case RoleRoot:
            return RoleGroupRoot
    }
    return RoleGuest
}

export function IsMember(role: number): boolean {
    return role == RoleRegularMember || role == RoleCertifiedMember
}

export function IsManager(role: number): boolean {
    return GetRoleGroup(role) == RoleGroupClassManager || GetRoleGroup(role) == RoleGroupGlobalManager || GetRoleGroup(role) == RoleGroupRoot
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
