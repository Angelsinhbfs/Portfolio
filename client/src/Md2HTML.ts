function Md2HTML(markdown: string): string {
    // Convert headers
    markdown = markdown.replace(/^### (.*$)/gim, '<h3>$1</h3>');
    markdown = markdown.replace(/^## (.*$)/gim, '<h2>$1</h2>');
    markdown = markdown.replace(/^# (.*$)/gim, '<h1>$1</h1>');

    // Convert bold text
    markdown = markdown.replace(/\*\*(.*)\*\*/gim, '<strong>$1</strong>');

    // Convert italic text
    markdown = markdown.replace(/\*(.*)\*/gim, '<em>$1</em>');

    // Convert unordered lists
    markdown = markdown.replace(/^\s*-\s+(.*$)/gim, '<ul><li>$1</li></ul>');

    // Convert links
    markdown = markdown.replace(/\[(.*?)\]\((.*?)\)/gim, '<a href="$2">$1</a>');

    // Convert images
    markdown = markdown.replace(/!\[(.*?)\]\((.*?)\)/gim, '<img alt="$1" src="$2" />');

    // Replace line breaks with <br>
    markdown = markdown.replace(/\n$/gim, '<br />');

    // Convert code blocks
    markdown = markdown.replace(/```(.*?)```/gim, '<code>$1</code>');

    // Convert unordered lists
    markdown = markdown.replace(/^\s*-\s+(.*$)/gim, '<ul><li>$1</li></ul>');

    // Convert ordered lists
    markdown = markdown.replace(/^\s*\d+\.\s+(.*$)/gim, '<ol><li>$1</li></ol>');

    // Convert checkboxes
    markdown = markdown.replace(/\[ \]/g, '<input type="checkbox" disabled />'); // Unchecked checkbox
    markdown = markdown.replace(/\[x\]/g, '<input type="checkbox" checked disabled />'); // Checked checkbox

    // Handle paragraphs
    markdown = markdown.replace(/(.+)(?:\r?\n)(?![\r\n])/g, '<p class="indented">$1</p>');

    return markdown.trim();
}
export default Md2HTML;