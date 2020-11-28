import {Component, Input, OnInit} from '@angular/core';
import {MessageGroup} from '../utils/messages-mapper';

@Component({
  selector: 'app-message-group',
  template: `
    <div class="d-flex flex-row mt-2 message-group">
      <div>
        <a class="sender" [routerLink]="'/users/' + messageGroup.sentById">
          <span
            class="badge badge-secondary user-badge">{{messageGroup.sentByUsername.substr(0, 1).toUpperCase()}}
          </span>
        </a>
      </div>

      <div class="flex-grow-1 message-list">
        <div class="group-header">
          <a class="sender text-dark" style="font-size: 0.85rem" [routerLink]="'/users/' + messageGroup.sentById">
            {{messageGroup.sentByUsername}}
          </a>&nbsp;
          <small class="sent-at text-muted" style="font-size: 0.75rem">
            {{messageGroup.formattedDate}} - {{messageGroup.messages[0].sentAtDate | date:"HH:MM aa"}}
          </small>
        </div>
        <div>
          <div class="message" *ngFor="let message of messageGroup.messages">
            <app-blocks [message]="message"></app-blocks>
          </div>
        </div>
      </div>

    </div>
  `,
  styles: [`
    .message-group {
      line-break: anywhere;
    }

    .group-header{
      line-height: 1rem;
    }
    .user-badge {
      width: 2rem;
      height: 2rem;
      padding-top: 0.65rem;
      margin-left: 1rem;
      position: relative;
      top: 0rem;
    }

    .message-list {
      font-size: 87%;
      margin-left: 0.5rem
    }
  `]
})
export class MessageGroupComponent {

  @Input()
  messageGroup: MessageGroup;

}
