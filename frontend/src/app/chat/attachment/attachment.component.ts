import {Component, Input} from '@angular/core';
import {Attachment} from '../../api/models';

@Component({
  selector: 'app-attachment',
  styles: [`
    .attachment-indicator {
      width: 0.25rem;
    }
  `],
  template: `
    <div class="attachment">
      <div class="d-flex flex-row">
        <div [style]="{'background': attachment.color}" class="attachment-indicator rounded mr-2"></div>
        <div>
          <app-block *ngFor="let block of attachment.blocks" [block]="block"></app-block>
        </div>
      </div>
    </div>`
})
export class AttachmentComponent {
  @Input()
  attachment: Attachment;
}
