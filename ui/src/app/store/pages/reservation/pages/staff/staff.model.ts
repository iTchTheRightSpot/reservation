export interface StaffModel {
  staff_id: string;
  name: string;
  image_key: string;
  bio: string;
}

export const DummyStaffModels: StaffModel[] = [
  {
    staff_id: '1',
    name: 'Tony',
    image_key: './salon-1.jpg',
    bio: 'Lorem ipsum dolor sit amet, consectetur adipisicing elit. Ea impedit maxime officiis rem unde.'
  },
  {
    staff_id: '2',
    name: 'benjamin',
    image_key: './salon-2.jpg',
    bio: 'Lorem ipsum dolor sit amet, consectetur adipisicing elit. Ab assumenda dolor error molestiae nobis qui quos recusandae sequi sit tempore?'
  },
  {
    staff_id: '3',
    name: 'phil',
    image_key: './salon-3.jpg',
    bio: 'Lorem ipsum dolor sit amet, consectetur adipisicing elit. Ea impedit maxime officiis rem unde.'
  }
];
