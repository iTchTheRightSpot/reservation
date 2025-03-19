export interface RolePermissionEntity {
  role: Role;
  permissions: Permission[];
}

export enum Role {
  STAFF = 'STAFF',
  DEVELOPER = 'DEVELOPER'
}

export enum Permission {
  READ = 'READ',
  WRITE = 'WRITE',
  DELETE = 'DELETE'
}

export interface CRMStaffModel {
  user_id: string;
  firstname: string;
  lastname: string;
  email: string;
  locked: boolean;
  image_key: string | null;
  bio: string;
  access_controls: RolePermissionEntity[];
}

export const DummyCRMStaffModels = (num: number) =>
  Array.from(
    { length: num },
    (_, index) =>
      ({
        user_id: `${index + 1}`,
        firstname: `Firstname-${index + 1}`,
        lastname: `Lastname-${index + 1}`,
        email: `email-${index + 1}@email.com`,
        locked: index % 2 == 0,
        image_key: index % 2 == 0 ? './salon-1.jpg' : null,
        bio: 'Lorem ipsum dolor sit amet, consectetur adipisicing elit. Dolores esse facilis impedit molestiae mollitia nesciunt nulla odit quod veritatis voluptate.',
        access_controls: [
          { role: Role.STAFF, permissions: [Permission.WRITE] },
          {
            role: Role.DEVELOPER,
            permissions: [Permission.WRITE, Permission.DELETE]
          }
        ]
      }) as CRMStaffModel
  );
