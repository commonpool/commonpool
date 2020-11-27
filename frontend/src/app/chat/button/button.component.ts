import {Component, Input, OnInit} from '@angular/core';
import {ButtonElement, SubmitAction} from '../../api/models';
import {BlocksService} from '../blocks.service';
import {BlockService} from '../block.service';

@Component({
  selector: 'app-button',
  templateUrl: './button.component.html',
  styleUrls: ['./button.component.css']
})
export class ButtonComponent implements OnInit {

  constructor(private blocksService: BlocksService, private blockService: BlockService) {
  }

  @Input()
  buttonElement: ButtonElement;

  ngOnInit(): void {
  }

  submit($event: MouseEvent) {
    console.log('OK');
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
