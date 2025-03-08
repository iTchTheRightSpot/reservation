import { BehaviorSubject, map } from 'rxjs';

export class Cache<K, V> {
  private readonly subject = new BehaviorSubject(new Map<K, V>());

  readonly getItem = (key: K) =>
    this.subject.asObservable().pipe(map(cache => cache.get(key)));

  readonly setItem = (key: K, value: V) =>
    this.subject.getValue().set(key, value);

  readonly clear = () =>
    this.subject.getValue().forEach((value, key, map) => map.delete(key));
}
