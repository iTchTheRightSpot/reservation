export interface StaffModel {
  staff_id: string;
  name: string;
  image_key: string | null;
  bio: string;
}

export const DummyStaffModels: StaffModel[] = [
  {
    staff_id: '1',
    name: 'Tony',
    image_key: './salon-1.jpg',
    bio: 'Ready to put a smile on your face🌞'
  },
  {
    staff_id: '2',
    name: 'benjamin',
    image_key: './salon-2.jpg',
    bio: 'Ready to put a smile on your face🌞'
  },
  {
    staff_id: '3',
    name: 'phil',
    image_key: './salon-3.jpg',
    bio: 'Ready to put a smile on your face🌞'
  }
];
