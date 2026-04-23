import { getUser, createUser } from './utils/helper';

interface User {
  id: string;
  name: string;
  email: string;
}

type UserRole = 'admin' | 'user' | 'guest';

/**
 * Fetches a user by their unique identifier.
 * @param id - The user's unique ID
 * @returns The user object or null if not found
 */
export async function fetchUser(id: string): Promise<User | null> {
  const user = await getUser(id);
  return user;
}

/**
 * Creates a new user with the given data.
 * @param data - The user creation payload
 */
export function addUser(data: Omit<User, 'id'>): User {
  return createUser(data);
}

export default { fetchUser, addUser };
