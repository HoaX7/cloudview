<script lang="ts">
  import Konva from "konva";
    import Group from "./Group.svelte";

    const path1 = `M56,46V12c0-0.553-0.447-1-1-1H9c-0.553,0-1,0.447-1,1v34c0,0.553,0.447,1,1,1h46
		C55.553,47,56,46.553,56,46z M54,45H10V13h44V45z`;
    // eslint-disable-next-line max-len
    const path2 = `M63,49h-3V11c0-2.211-1.789-4-4-4H8c-2.211,0-4,1.789-4,4v38H1c-0.553,0-1,0.447-1,1v3c0,2.211,1.789,4,4,4
		h56c2.211,0,4-1.789,4-4v-3C64,49.447,63.553,49,63,49z M6,11c0-1.104,0.896-2,2-2h48c1.104,0,2,0.896,2,2v38H36
		c-0.553,0-1,0.447-1,1v1h-6v-1c0-0.553-0.447-1-1-1H6V11z M62,53c0,1.104-0.896,2-2,2H4c-1.104,0-2-0.896-2-2v-2h25v
        1c0,0.553,0.447,1,1,1h8c0.553,0,1-0.447,1-1v-1h25V53z`;

    const Path = new Konva.Path({
    	data: path1,
    	x: 0,
    	y: 0,
    	scale: {
    		x: 2,
    		y: 2
    	},
    	fill: "#231F20",
    	draggable: false,
    	listening: false,
    });
    Path.cache();
  </script>
  
<Group config={{ y: 150 }} getHandler={(handle) => {
	handle.add(Path);
	const Path2El = Path.clone({
		listening: false,
		data: path2 
	});
	handle.add(Path2El);
}} />    