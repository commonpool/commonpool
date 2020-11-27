import {Component, Input, OnInit, ViewEncapsulation} from '@angular/core';
import {Attachment, Block, Message} from '../../api/models';
import {BlocksService} from '../blocks.service';

interface Payload {
  blocks: Block[];
  attachments: Attachment[];
}

@Component({
  selector: 'app-blocks',
  templateUrl: './blocks.component.html',
  styleUrls: ['./blocks.component.css'],
  encapsulation: ViewEncapsulation.None,
  providers: [BlocksService]
})
export class BlocksComponent implements OnInit {

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

  }

  ngOnInit(): void {
  }

}
