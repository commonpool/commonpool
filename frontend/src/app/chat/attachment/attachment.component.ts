import {Component, Input, OnInit} from '@angular/core';
import {Attachment} from '../../api/models';

@Component({
  selector: 'app-attachment',
  templateUrl: './attachment.component.html',
  styleUrls: ['./attachment.component.css']
})
export class AttachmentComponent implements OnInit {

  constructor() {
  }

  @Input()
  attachment: Attachment;

  ngOnInit(): void {
  }

}
