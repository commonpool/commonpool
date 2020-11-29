import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-sample',
  templateUrl: './sample.component.html',
  styleUrls: ['./sample.component.css']
})
export class SampleComponent implements OnInit {

  constructor() {
  }

  public data = `<commonpool-user id="ad1b6494-3814-4cfc-92ed-8af0e621137a"/>`;

  ngOnInit(): void {
  }

}
