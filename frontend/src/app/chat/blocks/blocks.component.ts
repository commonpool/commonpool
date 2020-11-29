import {Component, Input, OnInit, ViewEncapsulation} from '@angular/core';
import {Attachment, Block, Message} from '../../api/models';
import {BlocksService} from '../blocks.service';

@Component({
  selector: 'app-blocks',
  templateUrl: './blocks.component.html',
  styleUrls: ['./blocks.component.css'],
  encapsulation: ViewEncapsulation.None,
  providers: [BlocksService],
})
export class BlocksComponent {

  private _message: Message;
  @Input()
  public set message(value: Message) {
    this._message = value;
    this.svc.setMessage(value);
  }

  public get message(): Message {
    return this._message;
  }

  constructor(public svc: BlocksService) {
    console.log('>?');
  }

  trackBlock(index: number, block: Block): any {
    return index;
  }

}
