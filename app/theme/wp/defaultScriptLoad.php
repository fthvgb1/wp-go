<?php

/**
 * wordpress/wp-includes/script-loader.php:632
 * @var $scripts
 */
$con = 'package wp
import "github.com/fthvgb1/wp-go/safety"

func defaultScripts(m *safety.Map[string, *Script], suffix string){

';
foreach ($scripts->registered as $handle => $h) {
    $dep = 'nil';
    if ($h->deps) {
        $dep = '[]string{';
        $dep .= implode(',', array_map(fn($v) => '"' . $v . '"', $h->deps));
        $dep .= '}';
    }
    $con .= sprintf('m.Store("%s", NewScript("%s", "%s"+suffix+".js", %s, "%s", %s))
', $handle, $h->handle, str_replace('.min.js', '', $h->src), $dep, $h->ver, $h->args ?: 'nil');
}
$con .= '}';
file_put_contents('/tmp/scriptLoad.go', $con);


/**
 * put code to wordpress/wp-includes/class-wp-scripts.php:504
 * @param $handle
 * @param $object_name
 * @param $l10n
 * @return void
 */
function parseLocalize($handle, $object_name, $l10n): void
{
    if ('utils' != $handle) {
        $s = array_map('parseArr', array_keys($l10n), array_values($l10n));
        $m = implode("\n", $s);
        $x = sprintf('AddStaticLocalize("%s","%s",map[string]any{
		%s
	})
', $handle, $object_name, $m);
        file_put_contents('/tmp/bb.go', $x, FILE_APPEND);
    }
}

function parseArr($k, $v): string
{
    /**
     * @var $this object
     */
    if (is_array($v)) {
        if (array_diff_key(array_values($v), $v)) {
            $x = '';
            foreach ($v as $kk => $vv) {
                $x .= parseArr($kk, $vv);
            }
            return sprintf('"%s":map[string]any{
			%s
},', $k, $x);
        } else {
            $s = array_map(fn($ss) => sprintf('"%s"', $ss), $v);
            return sprintf('"%s":[]string{%s},', $k, implode(',', $s));
        }

    } else {
        return sprintf('"%s":`%s`,', $k, $v);
    }
}


/**
 * /var/www/html/wordpress/wp-includes/class-wp-theme-json.php:1712
 * @return void
 */
function presetsMetadata()
{
    $ss = <<<go
{
		path:            []string{%s},
		preventOverride: []string{%s},
		useDefaultNames: %s,
		valueKey:        "%s",
		valueFunc:       %s,
		cssVars:         "%s",
		classes: map[string]string{
			%s
		},
		properties: []string{%s},
	},
go;
    $s = '';
    foreach (static::PRESETS_METADATA as $val) {
        $arr = [];
        $arr[] = implode(',', array_map(fn($v) => '"' . $v . '"', $val['path']));
        $arr[] = (!$val['prevent_override']) ? '' : implode(',', array_map(fn($v) => '"' . $v . '"', $val['prevent_override']));
        $arr[] = $val['use_default_names'] ? 'true' : 'false';
        $arr[] = $val['value_key'] ?? '';
        $arr[] = $val['value_func'] ?? 'nil';
        $arr[] = $val['css_vars'];
        $arr[] = implode(",\n", array_map(fn($k, $v) => sprintf('"%s":"%s"', $k, $v), array_keys($val['classes']), array_values($val['classes'])));
        $arr[] = implode(',', array_map(fn($v) => '"' . $v . '"', $val['properties']));
        $s .= sprintf($ss, ...$arr);
    }
    echo $s;
}