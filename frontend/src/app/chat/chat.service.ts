import {Injectable} from '@angular/core';
import {ReplaySubject} from 'rxjs';

export enum ConversationType {
  Channel,
  Group
}

export interface SelectedConversation {
  type: ConversationType;
  id: string;
}

@Injectable({
  providedIn: 'root'
})
export class ChatService {

  private currentConversation = new ReplaySubject<SelectedConversation>();
  public currentConversation$ = this.currentConversation.asObservable();

  constructor() {
  }

  public setCurrentConversation(c: SelectedConversation) {
    this.currentConversation.next(c);
  }

}
