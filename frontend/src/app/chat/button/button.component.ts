import {Component, Input, OnInit} from '@angular/core';
import {ButtonElement, SubmitAction} from '../../api/models';
import {BlocksService} from '../blocks.service';
import {BlockService} from '../block.service';

@Component({
  selector: 'app-button',
  styles: [`
    button {
      font-size: 0.5rem;
      font-weight: bold;
      height: 1.75rem;
      padding-top: 0.135rem;
    }
  `],
  template: `
    <button class="btn btn-sm"
            [class.btn-outline-success]="buttonElement.style === 'primary'"
            [class.btn-outline-danger]="buttonElement.style==='danger'"
            [class.btn-outline-secondary]="buttonElement.style!=='danger' && buttonElement.style !== 'primary'"
            (click)="submit($event)"
    >
      <app-text-object [textObject]="buttonElement.text" [small]="true"></app-text-object>
    </button>
  `
})
export class ButtonComponent {

  constructor(private blocksService: BlocksService, private blockService: BlockService) {
  }

  @Input()
  buttonElement: ButtonElement;

  submit($event: MouseEvent) {
    this.blocksService.submitInteraction(new SubmitAction(
      this.blockService.getBlock().blockId,
      this.buttonElement.actionId,
      this.buttonElement.type,
      undefined,
      undefined,
      this.buttonElement.value
    ));
  }
}
