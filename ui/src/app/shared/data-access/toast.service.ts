import { Injectable } from '@angular/core';
import {
  concat,
  concatMap,
  Observable,
  of,
  Subject,
  switchMap,
  timer
} from 'rxjs';

export enum ToastEnum {
  NONE = 'NONE',
  ERROR = 'ERROR',
  SUCCESS = 'SUCCESS'
}

export interface IToast {
  state: ToastEnum;
  message: string;
}

@Injectable({
  providedIn: 'root'
})
export class ToastService {
  private readonly subject = new Subject<IToast>();

  readonly toast$: Observable<IToast> = this.subject.pipe(
    switchMap(obj =>
      concat(
        of(obj),
        timer(5000).pipe(
          concatMap(() => of({ state: ToastEnum.NONE, message: '' } as IToast))
        )
      )
    )
  );

  readonly message = (obj: IToast) => this.subject.next(obj);
}
